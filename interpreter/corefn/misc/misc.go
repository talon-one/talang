//go:generate go run ../generate_allop.go -pkg=misc

package misc

import (
	"errors"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var Misc = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name: "misc3",
		Arguments: []block.Kind{
			block.BlockKind,
		},
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		return args[0], nil
	},
}

var Noop = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:        "noop",
		Arguments:   []block.Kind{},
		Returns:     block.AnyKind,
		Description: "No operation",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		return nil, nil
	},
}

var ToString = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name: "toString",
		Arguments: []block.Kind{
			block.AnyKind,
		},
		Returns:     block.StringKind,
		Description: "Converts the parameter to a string",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		if len(args) < 1 {
			return nil, errors.New("invalid or missing arguments")
		}
		return block.NewString(args[0].String()), nil
	},
}
