//go:generate go run ../generate_allop.go -pkg=cmp

package cmp

import (
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	interpreter.RegisterCoreFunction(AllOperations()...)
}

var Equal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AtomKind,
			block.AtomKind,
			block.AtomKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the arguments are the same",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		for i := 1; i < len(args); i++ {
			if args[0].String != args[i].String {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var NotEqual = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "!=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AtomKind,
			block.AtomKind,
			block.AtomKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the arguments are not the same",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		for i := 1; i < len(args); i++ {
			if args[0].String == args[i].String {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThanDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ">",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is greather then the following",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Decimal
		for i := 1; i < len(args); i++ {
			if a.Cmp(args[i].Decimal) <= 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThanTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ">",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
			block.TimeKind,
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is greather then the following",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Time
		for i := 1; i < len(args); i++ {
			if !a.After(args[i].Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "<",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is less then the following",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Decimal
		for i := 1; i < len(args); i++ {
			if a.Cmp(args[i].Decimal) >= 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "<",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
			block.TimeKind,
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is less then the following",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Time
		for i := 1; i < len(args); i++ {
			if !a.Before(args[i].Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThanOrEqualDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ">=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is greather or equal then the following",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Decimal
		for i := 1; i < len(args); i++ {
			if a.Cmp(args[i].Decimal) < 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var GreaterThanOrEqualTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ">=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
			block.TimeKind,
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is greather or equal then the following",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Time
		for i := 0; i < len(args); i++ {
			if !a.Equal(args[i].Time) && !a.After(args[i].Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanOrEqualDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "<=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is less or equal then the following",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Decimal
		for i := 1; i < len(args); i++ {
			if a.Cmp(args[i].Decimal) > 0 {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var LessThanOrEqualTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "<=",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
			block.TimeKind,
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the first argument is less or equal then the following",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Time
		for i := 0; i < len(args); i++ {
			if !a.Equal(args[i].Time) && !a.Before(args[i].Time) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var BetweenDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "between",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
			block.DecimalKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the arguments are between the second last and the last argument",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)

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

var BetweenTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "between",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.TimeKind,
			block.TimeKind,
			block.TimeKind,
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Tests if the arguments are between the second last and the last argument",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)

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
