//go:generate go run ../generate_allop.go -pkg=misc

package misc

import (
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

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
			block.DecimalKind | block.StringKind | block.BoolKind | block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Converts the parameter to a string",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(args[0].String), nil
	},
}
