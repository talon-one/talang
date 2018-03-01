//go:generate go run ../generate_allop.go -pkg=math
package math

import (
	"github.com/ericlagergren/decimal"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	interpreter.RegisterCoreFunction(AllOperations()...)
}

var Add = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "+",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.DecimalKind,
		Description: "Adds the arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var d *decimal.Big
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d = d.Add(d, args[i].Decimal)
			}
		}
		return block.NewDecimal(d), nil
	},
}

var Sub = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "-",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.DecimalKind,
		Description: "Subtracts the arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var d *decimal.Big
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d = d.Sub(d, args[i].Decimal)
			}
		}
		return block.NewDecimal(d), nil
	},
}

var Mul = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "*",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.DecimalKind,
		Description: "Multiplies the arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var d *decimal.Big
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d = d.Mul(d, args[i].Decimal)
			}
		}
		return block.NewDecimal(d), nil
	},
}

var Div = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "/",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.DecimalKind,
		Description: "Divides the arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var d *decimal.Big
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d = d.Quo(d, args[i].Decimal)
			}
		}
		return block.NewDecimal(d), nil
	},
}

var Mod = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "mod",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.DecimalKind,
		Description: "Modulo the arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var d *decimal.Big
		for i := 0; i < len(args); i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d = d.Rem(d, args[i].Decimal)
			}
		}
		return block.NewDecimal(d), nil
	},
}

var Floor = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "floor",
		Arguments: []block.Kind{
			block.DecimalKind,
		},
		Returns:     block.DecimalKind,
		Description: "Floor the decimal argument",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		ctx := decimal.Context{Precision: args[0].Decimal.Context.Precision}
		if args[0].Decimal.Signbit() {
			ctx.RoundingMode = decimal.AwayFromZero
		} else {
			ctx.RoundingMode = decimal.ToZero
		}
		return block.NewDecimal(ctx.RoundToInt(args[0].Decimal)), nil
	},
}

var Ceil = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "ceil",
		Arguments: []block.Kind{
			block.DecimalKind,
		},
		Returns:     block.DecimalKind,
		Description: "Ceil the decimal argument",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		ctx := decimal.Context{Precision: args[0].Decimal.Context.Precision}
		if args[0].Decimal.Signbit() {
			ctx.RoundingMode = decimal.ToZero
		} else {
			ctx.RoundingMode = decimal.AwayFromZero
		}
		return block.NewDecimal(ctx.RoundToInt(args[0].Decimal)), nil
	},
}
