//go:generate go run ../generate_allop.go -pkg=math
package math

import (
	"github.com/talon-one/decimal"
	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/token"
)

func init() {
	if err := interpreter.RegisterCoreFunction(AllOperations()...); err != nil {
		panic(err)
	}
}

var Add = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "+",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Decimal,
		Description: "Adds the arguments",
		Example: `
(+ 1 1)                                                          ; returns 2
(+ 1 2 3)                                                        ; returns 6
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		var d decimal.Decimal
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d.Add(args[i].Decimal)
			}
		}
		return token.NewDecimal(d), nil
	},
}

var Sub = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "-",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Decimal,
		Description: "Subtracts the arguments",
		Example: `
(- 1 1)                                                          ; returns 0
(- 1 2 3)                                                        ; returns -4
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		var d decimal.Decimal
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d.Sub(args[i].Decimal)
			}
		}
		return token.NewDecimal(d), nil
	},
}

var Mul = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "*",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Decimal,
		Description: "Multiplies the arguments",
		Example: `
(* 1 2)                                                          ; returns 2
(* 1 2 3)                                                        ; returns 6
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		var d decimal.Decimal
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d.Mul(args[i].Decimal)
			}
		}
		return token.NewDecimal(d), nil
	},
}

var Div = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "/",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Decimal,
		Description: "Divides the arguments",
		Example: `
(/ 1 2)                                                          ; returns 0.5
(/ 1 2 3)                                                        ; returns 0.166666
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		var d decimal.Decimal
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d.Div(args[i].Decimal)
			}
		}
		return token.NewDecimal(d), nil
	},
}

var Mod = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "mod",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Decimal,
		Description: "Modulo the arguments",
		Example: `
(mod 1 2)                                                        ; returns 1
(mod 3 8 2)                                                      ; returns 1
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		var d decimal.Decimal
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d.Mod(args[i].Decimal)
			}
		}
		return token.NewDecimal(d), nil
	},
}

var Floor = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "floor",
		Arguments: []token.Kind{
			token.Decimal,
		},
		Returns:     token.Decimal,
		Description: "Floor the decimal argument",
		Example: `
(floor 2)                                                        ; returns 2
(floor 2.4)                                                      ; returns 2
(floor 2.5)                                                      ; returns 2
(floor 2.9)                                                      ; returns 2
(floor -2.7)                                                     ; returns -3
(floor -2)                                                       ; returns -2
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		args[0].Decimal.Floor()
		return token.NewDecimal(args[0].Decimal), nil
	},
}

var Ceil = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "ceil",
		Arguments: []token.Kind{
			token.Decimal,
		},
		Returns:     token.Decimal,
		Description: "Ceil the decimal argument",
		Example: `
(ceil 2)                                                         ; returns 2
(ceil 2.4)                                                       ; returns 3
(ceil 2.5)                                                       ; returns 3
(ceil 2.9)                                                       ; returns 3
(ceil -2.7)                                                      ; returns -2
(ceil -2)                                                        ; returns -2
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		args[0].Decimal.Ceil()
		return token.NewDecimal(args[0].Decimal), nil
	},
}
