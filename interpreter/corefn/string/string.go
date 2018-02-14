//go:generate go run ../generate_allop.go -pkg=string
package string

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var Add = shared.TaSignature{
	Name:       "+",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.StringKind,
	},
	Returns:     block.StringKind,
	Description: "Concat strings",
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		values := make([]string, argc)
		for i := 0; i < argc; i++ {
			values[i] = args[i].Text
		}
		return block.NewString(strings.Join(values, "")), nil
	},
}

var Concat = shared.TaSignature{
	Name:        "concat",
	IsVariadic:  Add.IsVariadic,
	Arguments:   Add.Arguments,
	Returns:     Add.Returns,
	Description: Add.Description,
	Func:        Add.Func,
}

var Contains = shared.TaSignature{
	Name:       "contains",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.StringKind,
		block.StringKind,
	},
	Returns:     block.BoolKind,
	Description: "Returns wether the first argument exists in the following arguments",
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}

		for i := 1; i < argc; i++ {
			if !strings.Contains(args[0].Text, args[i].Text) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var NotContains = shared.TaSignature{
	Name:       "notContains",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.StringKind,
		block.StringKind,
	},
	Returns:     block.BoolKind,
	Description: "Returns wether the first argument does not exist in the following arguments",
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}

		for i := 1; i < argc; i++ {
			if strings.Contains(args[0].Text, args[i].Text) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var StartsWith = shared.TaSignature{
	Name:       "startsWith",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.StringKind,
		block.StringKind,
	},
	Returns:     block.BoolKind,
	Description: "Returns wether the first argument is the prefix of the following arguments",
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}

		for i := 1; i < argc; i++ {
			if !strings.HasPrefix(args[0].Text, args[i].Text) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}
var EndsWith = shared.TaSignature{
	Name:       "endsWith",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.StringKind,
		block.StringKind,
	},
	Returns:     block.BoolKind,
	Description: "Returns wether the first argument is the suffix of the following arguments",
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}

		for i := 1; i < argc; i++ {
			if !strings.HasSuffix(args[0].Text, args[i].Text) {
				return block.NewBool(false), nil
			}
		}
		return block.NewBool(true), nil
	},
}

var Regexp = shared.TaSignature{
	Name:       "~",
	IsVariadic: false,
	Arguments: []block.Kind{
		block.StringKind,
		block.StringKind,
	},
	Returns:     block.BoolKind,
	Description: "Returns wether the first argument matches the regular expression in the second argument",
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}
		re, err := regexp.Compile(args[1].String())
		if err != nil {
			return block.NewBool(false), err
		}
		return block.NewBool(re.MatchString(args[0].String())), nil
	},
}
