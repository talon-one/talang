//go:generate go run generate.go -pkg=string
package string

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var Contains = shared.TaSignature{
	Name:       "contains",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.StringKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}

		for i := 1; i < argc; i++ {
			if !strings.Contains(args[0].Text, args[i].Text) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var NotContains = shared.TaSignature{
	Name:       "notContains",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.StringKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}

		for i := 1; i < argc; i++ {
			if strings.Contains(args[0].Text, args[i].Text) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}
