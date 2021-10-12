//go:generate go run ../generate_allop.go -pkg=misc

package misc

import (
	"github.com/pkg/errors"

	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/token"
)

func init() {
	if err := interpreter.RegisterCoreFunction(AllOperations()...); err != nil {
		panic(err)
	}
}

var Noop = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:        "noop",
		Arguments:   []token.Kind{},
		Returns:     token.Any,
		Description: "No operation",
		Example:     `(noop)`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		return nil, nil
	},
}

var ToString = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "toString",
		Arguments: []token.Kind{
			token.Decimal | token.String | token.Boolean | token.Time,
		},
		Returns:     token.String,
		Description: "Converts the parameter to a string",
		Example: `
(toString 1)                                                     ; returns "1"
(toString true)                                                  ; returns "true"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		return token.NewString(args[0].String), nil
	},
}

var Not = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "not",
		Arguments: []token.Kind{
			token.Boolean,
		},
		Returns:     token.Boolean,
		Description: "Inverts the argument",
		Example: `
(not false)                                                      ; returns "true"
(not (not false))                                                ; returns "false"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		return token.NewBool(!args[0].Bool), nil
	},
}

var Catch = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "catch",
		Arguments: []token.Kind{
			token.Any,
			token.Any,
		},
		Returns:     token.Any,
		Description: "Evaluate & return the second argument. If any errors occur, return the first argument instead",
		Example: `
catch "Edward" (. Profile Name)                                  ; returns "Edward"
catch 22 (. Profile Age)                                         ; returns 46
catch 22 2                                                       ; returns 22
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		scope := interp.NewScope()
		err := scope.Evaluate(args[1])
		if err != nil {
			return args[0], nil
		}
		return args[1], nil
	},
}

var Do = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "do",
		Arguments: []token.Kind{
			token.Atom | token.Collection,
			token.String,
			token.Token,
		},
		Returns:     token.Any,
		Description: "Apply a block to a value",
		Example: `
do (list 1 2 3) Item (. Item))                                   ; returns 1 2 3
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		value := args[0]
		bindingName := args[1].String
		blockToRun := args[2]

		scope := interp.NewScope()
		scope.Set(bindingName, value)

		if err := scope.Evaluate(blockToRun); err != nil {
			return nil, err
		}
		return blockToRun, nil
	},
}

var DoLegacy = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "do",
		Arguments: []token.Kind{
			token.Atom | token.Collection,
			token.Token,
		},
		Returns:     token.Any,
		Description: "Apply a block to a value",
		Example: `
do (list 1 2 3) ((Item) (. Item)))                               ; returns 1 2 3
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		if len(args[1].Children) == 2 && args[1].Children[0].IsBlock() {
			return Do.Func(interp, args[0], args[1].Children[0], args[1].Children[1])
		}
		return nil, errors.New("Missing or invalid binding")
	},
}

var SafeRead = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ".|",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Any,
			token.Token, // binding
		},
		Returns:     token.Any,
		Description: "Safe read a binding",
		Example: `
.| boo (. List)                                                  ; returns "XJK_992" (assuming $List = "XJK_992")
.| boo (. Meh)                                                   ; returns "boo" (assuming $Meh is not set)
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		blockToRun := args[1]
		fallback := args[0]
		if err := interp.Evaluate(blockToRun); err != nil {
			return fallback, nil
		}
		if blockToRun.Kind != fallback.Kind {
			return nil, errors.Errorf("Cannot reconciliate type %s with %s", fallback.Kind.String(), blockToRun.Kind.String())
		}
		return blockToRun, nil
	},
}

var Identity = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "identity",
		Arguments: []token.Kind{token.Any},
		Returns:     token.Any,
		Description: "Returns its argument",
		Example: `
(identity 1)                                                     ; returns 1
(identity true)                                                  ; returns true
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		return args[0], nil
	},
}
