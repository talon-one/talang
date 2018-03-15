//go:generate go run ../generate_allop.go -pkg=misc

package misc

import (
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	if err := interpreter.RegisterCoreFunction(AllOperations()...); err != nil {
		panic(err)
	}
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
(toString 1)                                                     ; returns "1"
(toString true)                                                  ; returns "true"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(args[0].String), nil
	},
}

var Not = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "not",
		Arguments: []block.Kind{
			block.BoolKind,
		},
		Returns:     block.BoolKind,
		Description: "Inverts the argument",
		Example: `
(not false)                                                      ; returns "true"
(not (not false))                                                ; returns "false"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewBool(!args[0].Bool), nil
	},
}

var Catch = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "catch",
		Arguments: []block.Kind{
			block.AnyKind,
			block.AnyKind,
		},
		Returns:     block.AnyKind,
		Description: "Evaluate & return the second argument. If any errors occur, return the first argument instead",
		Example: `
catch "Edward" (. Profile Name)                                  ; returns "Edward"
catch 22 (. Profile Age)                                         ; returns 46
catch 22 2                                                       ; returns 22
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		scope := interp.NewScope()
		err := scope.Evaluate(args[1])
		if err != nil {
			return args[0], nil
		}
		return args[1], nil
	},
}
