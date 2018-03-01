//go:generate go run ../generate_allop.go -pkg=string
package string

import (
	"regexp"
	"strings"

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
			block.StringKind,
			block.StringKind,
			block.StringKind,
		},
		Returns:     block.StringKind,
		Description: "Concat strings",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		values := make([]string, argc)
		for i := 0; i < argc; i++ {
			values[i] = args[i].String
		}
		return block.NewString(strings.Join(values, "")), nil
	},
}

var Concat = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:        "concat",
		IsVariadic:  Add.IsVariadic,
		Arguments:   Add.Arguments,
		Returns:     Add.Returns,
		Description: Add.Description,
	},
	Func: Add.Func,
}

var Contains = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "contains",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.StringKind,
			block.StringKind,
			block.StringKind,
		},
		Returns:     block.BoolKind,
		Description: "Returns wether the first argument exists in the following arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		for i := 1; i < len(args); i++ {
			if !strings.Contains(args[0].String, args[i].String) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var NotContains = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "notContains",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.StringKind,
			block.StringKind,
			block.StringKind,
		},
		Returns:     block.BoolKind,
		Description: "Returns wether the first argument does not exist in the following arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		for i := 1; i < len(args); i++ {
			if strings.Contains(args[0].String, args[i].String) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var StartsWith = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "startsWith",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.StringKind,
			block.StringKind,
			block.StringKind,
		},
		Returns:     block.BoolKind,
		Description: "Returns wether the first argument is the prefix of the following arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		for i := 1; i < len(args); i++ {
			if !strings.HasPrefix(args[0].String, args[i].String) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}
var EndsWith = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "endsWith",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.StringKind,
			block.StringKind,
			block.StringKind,
		},
		Returns:     block.BoolKind,
		Description: "Returns wether the first argument is the suffix of the following arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		for i := 1; i < len(args); i++ {
			if !strings.HasSuffix(args[0].String, args[i].String) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var Regexp = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "~",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.StringKind,
			block.StringKind,
			block.StringKind,
		},
		Returns:     block.BoolKind,
		Description: "Returns wether the first argument (regex) matches all of the following arguments",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		re, err := regexp.Compile(args[0].String)
		if err != nil {
			return block.NewBool(false), err
		}

		for i := 1; i < len(args); i++ {
			if !re.MatchString(args[i].String) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}
