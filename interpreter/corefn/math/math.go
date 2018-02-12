//go:generate go run ../generate_allop.go -pkg=math
package math

import (
	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var Add = shared.TaSignature{
	Name:       "+",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Returns:     block.DecimalKind,
	Description: "Adds the arguments",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}
		var d *decimal.Big
		for i := 0; i < argc; i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				d = d.Add(d, args[i].Decimal)
			}
		}
		return block.NewDecimal(d), nil
	},
}

var Sub = shared.TaSignature{
	Name:       "-",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Returns:     block.DecimalKind,
	Description: "Substracts the arguments",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		if len(args) < 2 {
			return nil, errors.New("invalid or missing arguments")
		}
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

var Mul = shared.TaSignature{
	Name:       "*",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Returns:     block.DecimalKind,
	Description: "Multiplies the arguments",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		if len(args) < 2 {
			return nil, errors.New("invalid or missing arguments")
		}
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

var Div = shared.TaSignature{
	Name:       "/",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Returns:     block.DecimalKind,
	Description: "Divides the arguments",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		if len(args) < 2 {
			return nil, errors.New("invalid or missing arguments")
		}
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

var Mod = shared.TaSignature{
	Name:       "mod",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Returns:     block.DecimalKind,
	Description: "Modulo the arguments",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		if len(args) < 2 {
			return nil, errors.New("invalid or missing arguments")
		}
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

var Floor = shared.TaSignature{
	Name: "floor",
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Returns:     block.DecimalKind,
	Description: "Floor the decimal argument",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		if len(args) != 1 {
			return nil, errors.New("invalid or missing arguments")
		}
		ctx := decimal.Context{Precision: args[0].Decimal.Context.Precision}
		if args[0].Decimal.Signbit() {
			ctx.RoundingMode = decimal.AwayFromZero
		} else {
			ctx.RoundingMode = decimal.ToZero
		}
		return block.NewDecimal(ctx.RoundToInt(args[0].Decimal)), nil
	},
}

var Ceil = shared.TaSignature{
	Name: "ceil",
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Returns:     block.DecimalKind,
	Description: "Ceil the decimal argument",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		if len(args) != 1 {
			return nil, errors.New("invalid or missing arguments")
		}

		ctx := decimal.Context{Precision: args[0].Decimal.Context.Precision}
		if args[0].Decimal.Signbit() {
			ctx.RoundingMode = decimal.ToZero
		} else {
			ctx.RoundingMode = decimal.AwayFromZero
		}
		return block.NewDecimal(ctx.RoundToInt(args[0].Decimal)), nil
	},
}
