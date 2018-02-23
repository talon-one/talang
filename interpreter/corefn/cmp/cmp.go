//go:generate go run ../generate_allop.go -pkg=cmp

package cmp

import (
	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var Equal = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AnyKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the arguments are the same",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
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

var NotEqual = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "!=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AnyKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the arguments are not the same",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
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

var GreaterThanDecimal = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       ">",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is greather then the following",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		a := args[0].Decimal
		for i := 1; i < argc; i++ {
			if a.Cmp(args[i].Decimal) <= 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThanTime = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       ">",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is greather then the following",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		a := args[0].Time
		for i := 1; i < argc; i++ {
			if !a.After(args[i].Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanDecimal = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "<",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is less then the following",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		a := args[0].Decimal
		for i := 1; i < argc; i++ {
			if a.Cmp(args[i].Decimal) >= 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanTime = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "<",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is less then the following",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		a := args[0].Time
		for i := 1; i < argc; i++ {
			if !a.Before(args[i].Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThanOrEqualDecimal = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       ">=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is greather or equal then the following",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		a := args[0].Decimal
		for i := 1; i < argc; i++ {
			if a.Cmp(args[i].Decimal) < 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThanOrEqualTime = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       ">=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is greather or equal then the following",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		a := args[0].Time
		for i := 0; i < argc; i++ {
			if !a.Equal(args[i].Time) && !a.After(args[i].Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanOrEqualDecimal = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "<=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is less or equal then the following",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		a := args[0].Decimal
		for i := 1; i < argc; i++ {
			if a.Cmp(args[i].Decimal) > 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanOrEqualTime = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "<=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is less or equal then the following",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		a := args[0].Time
		for i := 0; i < argc; i++ {
			if !a.Equal(args[i].Time) && !a.Before(args[i].Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var BetweenDecimal = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "between",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the arguments are between the second last and the last argument",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
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

var BetweenTime = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "between",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the arguments are between the second last and the last argument",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 3 {
			return block.NewBool(false), errors.New("invalid or missing arguments")
		}

		min := args[argc-2]
		max := args[argc-1]

		argc -= 2

		for i := 0; i < argc; i++ {
			if args[i].Time.Before(min.Time) || args[i].Time.After(max.Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}
