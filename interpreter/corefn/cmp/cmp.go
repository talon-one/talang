//go:generate go run generate.go -pkg=cmp
package cmp

import (
	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var Equal = shared.TaSignature{
	Name:       "=",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.AnyKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}

		for i := 1; i < argc; i++ {
			if args[0].Text != args[i].Text {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var NotEqual = shared.TaSignature{
	Name:       "!=",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.AnyKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		for i := 1; i < argc; i++ {
			if args[0].Text == args[i].Text {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThan = shared.TaSignature{
	Name:       ">",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		var d *decimal.Big
		for i := 0; i < argc; i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				if d.Cmp(args[i].Decimal) <= 0 {
					return block.NewBool(false), nil
				}
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThan = shared.TaSignature{
	Name:       "<",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		var d *decimal.Big
		for i := 0; i < argc; i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				if d.Cmp(args[i].Decimal) >= 0 {
					return block.NewBool(false), nil
				}
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThanOrEqual = shared.TaSignature{
	Name:       ">=",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		var d *decimal.Big
		for i := 0; i < argc; i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				if d.Cmp(args[i].Decimal) < 0 {
					return block.NewBool(false), nil
				}
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanOrEqual = shared.TaSignature{
	Name:       "<=",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		var d *decimal.Big
		for i := 0; i < argc; i++ {
			if i == 0 {
				d = args[i].Decimal
			} else {
				if d.Cmp(args[i].Decimal) > 0 {
					return block.NewBool(false), nil
				}
			}
		}
		return block.NewBool(true), nil
	},
}

var Between = shared.TaSignature{
	Name:       "between",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.DecimalKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 3 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		min := args[argc-2]
		max := args[argc-1]

		argc -= 2

		for i := 0; i < argc; i++ {
			if args[i].Decimal.Cmp(min.Decimal) < 0 || args[i].Decimal.Cmp(max.Decimal) > 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}
