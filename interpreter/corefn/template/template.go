//go:generate go run ../generate_allop.go -pkg=template
package template

import (
	"context"

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
	},
	Returns:     block.BlockKind,
	Description: "Resolve a template",
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		m := getMap(interp)
		if b, ok := m[args[0].Text]; ok {
			return &b, nil
		}
		return nil, errors.Errorf("template `%s' not found", args[0].Text)
	},
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
		m[args[0].Text] = *args[1]
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
