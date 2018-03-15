//go:generate go run ../generate_allop.go -pkg=string
package string

import (
	"regexp"
	"strings"

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
			token.String,
			token.String,
			token.String,
		},
		Returns:     token.String,
		Description: "Concat strings",
		Example: `
(+ "Hello" " " "World")                                          ; returns "Hello World"
(+ "Hello" " " (toString (+ 1 2)))                               ; returns "Hello 3"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		argc := len(args)
		values := make([]string, argc)
		for i := 0; i < argc; i++ {
			values[i] = args[i].String
		}
		return token.NewString(strings.Join(values, "")), nil
	},
}

var Concat = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:        "concat",
		IsVariadic:  Add.IsVariadic,
		Arguments:   Add.Arguments,
		Returns:     Add.Returns,
		Description: Add.Description,
		Example:     Add.Example,
	},
	Func: Add.Func,
}

var Contains = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "contains",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.String,
			token.String,
			token.String,
		},
		Returns:     token.Bool,
		Description: "Returns wether the first argument exists in the following arguments",
		Example: `
(contains "Hello" "Hello World")                                 ; returns true
(contains "Hello" "World")                                       ; returns false
(contains "Hello" "Hello World" "Hello Universe")                ; returns true
(contains "World" "Hello World" "Hello Universe")                ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		for i := 1; i < len(args); i++ {
			if !strings.Contains(args[0].String, args[i].String) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var NotContains = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "notContains",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.String,
			token.String,
			token.String,
		},
		Returns:     token.Bool,
		Description: "Returns wether the first argument does not exist in the following arguments",
		Example: `
(notContains "Hello" "Hello World")                              ; returns false
(notContains "Hello" "World")                                    ; returns true
(notContains "Hello" "Hello World" "Hello Universe")             ; returns false
(notContains "World" "Hello World" "Hello Universe")             ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		for i := 1; i < len(args); i++ {
			if strings.Contains(args[0].String, args[i].String) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var StartsWith = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "startsWith",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.String,
			token.String,
			token.String,
		},
		Returns:     token.Bool,
		Description: "Returns wether the first argument is the prefix of the following arguments",
		Example: `
(startsWith "Hello" "Hello World")                               ; returns true
(startsWith "Hello" "World")                                     ; returns false
(startsWith "Hello" "Hello World" "Hello Universe")              ; returns true
(startsWith "Hello" "Hello World" "Hell Universe")               ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		for i := 1; i < len(args); i++ {
			if !strings.HasPrefix(args[0].String, args[i].String) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}
var EndsWith = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "endsWith",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.String,
			token.String,
			token.String,
		},
		Returns:     token.Bool,
		Description: "Returns wether the first argument is the suffix of the following arguments",
		Example: `
(endsWith "World" "Hello World")                                 ; returns true
(endsWith "World" "Hello Universe")                              ; returns false
(endsWith "World" "Hello World" "Hello Universe")                ; returns false
(endsWith "World" "Hello World" "By World")                      ; returns true
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		for i := 1; i < len(args); i++ {
			if !strings.HasSuffix(args[0].String, args[i].String) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var Regexp = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "~",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.String,
			token.String,
			token.String,
		},
		Returns:     token.Bool,
		Description: "Returns wether the first argument (regex) matches all of the following arguments",
		Example: `
(~ "[a-z\s]*" "Hello World")                                     ; returns true
(~ "[a-z\s]*" "Hello W0rld")                                     ; returns false
(~ "[a-z\s]*" "Hello World" "Hello Universe")                    ; returns true
(~ "[a-z\s]*" "Hello W0rld" "Hello Universe")                    ; returns false
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		re, err := regexp.Compile(args[0].String)
		if err != nil {
			return token.NewBool(false), err
		}

		for i := 1; i < len(args); i++ {
			if !re.MatchString(args[i].String) {
				return token.NewBool(false), nil
			}
		}
		return token.NewBool(true), nil
	},
}

var LastName = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "lastName",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.String,
		},
		Returns:     token.String,
		Description: "Extract the last word (space-separated) from a string",
		Example: `
(lastName "Alex Unger")                                          ; returns "Unger"
(lastName "Mr Foo Bar")                                          ; returns "Bar"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		words := strings.Split(args[0].String, " ")
		return token.NewString(words[len(words)-1]), nil
	},
}

var FirstName = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "firstName",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.String,
		},
		Returns:     token.String,
		Description: "Extract all but the last word (space-separated) from a string",
		Example: `
(firstName "Alex Unger")                                         ; returns "Alex"
(firstName "Mr Foo Bar")                                         ; returns "Mr"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		words := strings.Split(args[0].String, " ")
		return token.NewString(words[0]), nil
	},
}
