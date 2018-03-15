//go:generate go run ../generate_allop.go -pkg=cmp

package cmp

import (
	"github.com/talon-one/talang/token"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	if err := interpreter.RegisterCoreFunction(AllOperations()...); err != nil {
		panic(err)
	}
}

var Equal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "=",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Atom,
			token.Atom,
			token.Atom,
		},
		Returns:     token.Bool,
		Description: "Tests if the arguments are the same",
		Example: `
(= 1 1)                                                          ; compares decimals, returns true
(= "Hello World" "Hello World")                                  ; compares strings, returns true
(= true true)                                                    ; compares booleans, returns true
(= 2006-01-02T15:04:05Z 2006-01-02T15:04:05Z)                    ; compares time, returns true
(= 1 "1")                                                        ; returns true
(= "Hello" "Bye")                                                ; returns false
(= "Hello" "Hello" "Bye")                                        ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		for i := 1; i < len(args); i++ {
			if args[0].String != args[i].String {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var NotEqual = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "!=",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Atom,
			token.Atom,
			token.Atom,
		},
		Returns:     token.Bool,
		Description: "Tests if the arguments are not the same",
		Example: `
(!= 1 1)                                                         ; compares decimals, returns false
(!= "Hello World" "Hello World")                                 ; compares strings, returns false
(!= true true)                                                   ; compares booleans, returns false
(!= 2006-01-02T15:04:05Z 2006-01-02T15:04:05Z)                   ; compares time, returns false
(!= 1 "1")                                                       ; returns false
(!= "Hello" "Bye")                                               ; returns true
(!= "Hello" "Hello" "Bye")                                       ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		for i := 1; i < len(args); i++ {
			if args[0].String == args[i].String {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var GreaterThanDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ">",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Bool,
		Description: "Tests if the first argument is greather then the following",
		Example: `
(> 0 1)                                                          ; returns false
(> 1 1)                                                          ; returns false
(> 2 1)                                                          ; returns true
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Decimal
		for i := 1; i < len(args); i++ {
			if a.Cmp(args[i].Decimal) <= 0 {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var GreaterThanTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ">",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Time,
			token.Time,
			token.Time,
		},
		Returns:     token.Bool,
		Description: "Tests if the first argument is greather then the following",
		Example: `
(> 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns false
(> 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns false
(> 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns true
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Time
		for i := 1; i < len(args); i++ {
			if !a.After(args[i].Time) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var LessThanDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "<",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Bool,
		Description: "Tests if the first argument is less then the following",
		Example: `
(< 0 1)                                                          ; returns true
(< 1 1)                                                          ; returns false
(< 2 1)                                                          ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Decimal
		for i := 1; i < len(args); i++ {
			if a.Cmp(args[i].Decimal) >= 0 {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var LessThanTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "<",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Time,
			token.Time,
			token.Time,
		},
		Returns:     token.Bool,
		Description: "Tests if the first argument is less then the following",
		Example: `
(< 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns true
(< 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns false
(< 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)                    ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Time
		for i := 1; i < len(args); i++ {
			if !a.Before(args[i].Time) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var GreaterThanOrEqualDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ">=",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Bool,
		Description: "Tests if the first argument is greather or equal then the following",
		Example: `
(>= 0 1)                                                         ; returns false
(>= 1 1)                                                         ; returns true
(>= 2 1)                                                         ; returns true
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Decimal
		for i := 1; i < len(args); i++ {
			if a.Cmp(args[i].Decimal) < 0 {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var GreaterThanOrEqualTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       ">=",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Time,
			token.Time,
			token.Time,
		},
		Returns:     token.Bool,
		Description: "Tests if the first argument is greather or equal then the following",
		Example: `
(>= 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns false
(>= 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns true
(>= 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns true
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Time
		for i := 0; i < len(args); i++ {
			if !a.Equal(args[i].Time) && !a.After(args[i].Time) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var LessThanOrEqualDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "<=",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Bool,
		Description: "Tests if the first argument is less or equal then the following",
		Example: `
(<= 0 1)                                                         ; returns true
(<= 1 1)                                                         ; returns true
(<= 2 1)                                                         ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Decimal
		for i := 1; i < len(args); i++ {
			if a.Cmp(args[i].Decimal) > 0 {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var LessThanOrEqualTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "<=",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Time,
			token.Time,
			token.Time,
		},
		Returns:     token.Bool,
		Description: "Tests if the first argument is less or equal then the following",
		Example: `
(<= 2006-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns true
(<= 2007-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns true
(<= 2008-01-02T15:04:05Z 2007-01-02T15:04:05Z)                   ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Time
		for i := 0; i < len(args); i++ {
			if !a.Equal(args[i].Time) && !a.Before(args[i].Time) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var BetweenDecimal = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "between",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Decimal,
			token.Decimal,
			token.Decimal,
			token.Decimal,
		},
		Returns:     token.Bool,
		Description: "Tests if the arguments are between the second last and the last argument",
		Example: `
(between 1 0 3)                                                  ; returns true, (1 is between 0 and 3)
(between 1 2 0 3)                                                ; returns true, (1 and 2 are between 0 and 3)
(between 0 0 2)                                                  ; returns false
(between 2 0 2)                                                  ; returns false
(between 1 4 0 3)                                                ; returns false, (1 is between 0 and 3, 4 is not)
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		argc := len(args)

		min := args[argc-2]
		max := args[argc-1]

		argc -= 2

		for i := 0; i < argc; i++ {
			if args[i].Decimal.Cmp(min.Decimal) <= 0 || args[i].Decimal.Cmp(max.Decimal) >= 0 {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var BetweenTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "between",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Time,
			token.Time,
			token.Time,
			token.Time,
		},
		Returns:     token.Bool,
		Description: "Tests if the arguments are between the second last and the last argument",
		Example: `
(between 2007-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)                        ; returns true, (2007-01-02T00:00:00Z is between 2006-01-02T00:00:00Z and 3)
(between 2007-01-02T00:00:00Z 2008-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)   ; returns true, (2007-01-02T00:00:00Z and 2008-01-02T00:00:00Z are between 2006-01-02T00:00:00Z and 2009-01-02T00:00:00Z)
(between 2006-01-02T00:00:00Z 2006-01-02T00:00:00Z 2008-01-02T00:00:00Z)                        ; returns false
(between 2008-01-02T00:00:00Z 2006-01-02T00:00:00Z 2008-01-02T00:00:00Z)                        ; returns false
(between 2007-01-02T00:00:00Z 2010-01-02T00:00:00Z 2006-01-02T00:00:00Z 2009-01-02T00:00:00Z)   ; returns false, (2007-01-02T00:00:00Z is between 2006-01-02T00:00:00Z and 2009-01-02T00:00:00Z, 2010-01-02T00:00:00Z is not)
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		argc := len(args)

		min := args[argc-2]
		max := args[argc-1]

		argc -= 2

		for i := 0; i < argc; i++ {
			if args[i].Time.Equal(min.Time) || args[i].Time.Equal(max.Time) || args[i].Time.Before(min.Time) || args[i].Time.After(max.Time) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}
