//go:generate go run ../generate_allop.go -pkg=template
package template

import (
	"context"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

type contextKey int

const (
	templateKey contextKey = iota
)

var GetTemplate = shared.TaSignature{
	Name:       "!",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.StringKind,
		block.AnyKind,
	},
	Returns:     block.BlockKind,
	Description: "Resolve a template",
	Func:        getTemplateFunc,
}

func getTemplateFunc(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
	argc := len(args)
	if argc < 1 {
		return nil, errors.New("invalid or missing arguments")
	}
	m := getMap(interp)
	if b, ok := m[strings.ToLower(args[0].Text)]; ok {
		// the template has arguments
		// replace all variables in the block
		if argc > 1 {
			if _, err := replaceVariables(&b, args[1:]...); err != nil {
				return nil, err
			}
		}
		return &b, nil
	}
	if interp.Parent != nil {
		return getTemplateFunc(interp.Parent, args...)
	}
	return nil, errors.Errorf("template `%s' not found", args[0].Text)
}

func replaceVariables(b *block.Block, args ...*block.Block) (int, error) {
	total := 0

	var replaced int

replace:
	replaced = 0
	for i := 0; i < len(args); i++ {
		replaced += replaceVariable(b, strconv.Itoa(i), args[i])
	}

	total += replaced

	if replaced > 0 {
		goto replace
	}
	return total, nil
}

func replaceVariable(source *block.Block, name string, replace *block.Block) (replaced int) {
	if len(source.Children) <= 0 {
		return replaced
	}

	if source.Text == "#" {
		if strings.EqualFold(source.Children[0].Text, name) {
			*source = *replace
			replaced++
		}
	}

	for i := 0; i < len(source.Children); i++ {
		replaced += replaceVariable(source.Children[i], name, replace)
	}

	return replaced
}

var SetTemplate = shared.TaSignature{
	Name: "setTemplate",
	Arguments: []block.Kind{
		block.StringKind,
		block.BlockKind,
	},
	Returns:     block.BlockKind,
	Description: "Set a template",
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		if len(args) < 2 {
			return nil, errors.New("invalid or missing arguments")
		}
		m := getMap(interp)
		m[strings.ToLower(args[0].Text)] = *args[1]
		return nil, nil
	},
}

func getMap(interp *shared.Interpreter) map[string]block.Block {
	if interp.Context != nil {
		if m := interp.Context.Value(templateKey); m != nil {
			if mp, ok := m.(map[string]block.Block); ok {
				return mp
			}
		}
	}
	mp := make(map[string]block.Block)
	interp.Context = context.WithValue(interp.Context, templateKey, mp)
	return mp
}

func Set(interp *shared.Interpreter, name string, block block.Block) error {
	m := getMap(interp)
	m[strings.ToLower(name)] = block
	return nil
}
