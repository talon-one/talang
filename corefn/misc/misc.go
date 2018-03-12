//go:generate go run ../generate_allop.go -pkg=misc

package misc

import (
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	interpreter.RegisterCoreFunction(AllOperations()...)
}

var Noop = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:        "noop",
		Arguments:   []block.Kind{},
		Returns:     block.AnyKind,
		Description: "No operation",
		Example:     `(noop)`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return nil, nil
	},
}

var ToString = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "toString",
		Arguments: []block.Kind{
			block.DecimalKind | block.StringKind | block.BoolKind | block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Converts the parameter to a string",
		Example: `
(toString 1)                                                      // returns "1"
(toString true)                                                   // returns "true"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(args[0].String), nil
	},
}
