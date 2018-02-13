//go:generate go run ../generate_allop.go -pkg=list
package list

import (
	"github.com/pkg/errors"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var List = shared.TaSignature{
	Name:       "list",
	IsVariadic: true,
	Arguments: []block.Kind{
		block.AtomKind,
	},
	Returns:     block.BlockKind,
	Description: "Create a list out of the children",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		return block.New("", args...), nil
	},
}

var Head = shared.TaSignature{
	Name:       "head",
	IsVariadic: false,
	Arguments: []block.Kind{
		block.BlockKind,
	},
	Returns:     block.BlockKind,
	Description: "Returns the first item in the list",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		if len(args) < 1 {
			return nil, errors.New("invalid or missing arguments")
		}
		return args[0], nil
	},
}

var Tail = shared.TaSignature{
	Name:       "tail",
	IsVariadic: false,
	Arguments: []block.Kind{
		block.BlockKind,
	},
	Returns:     block.BlockKind,
	Description: "Returns list without the first item",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		if len(args) < 1 {
			return nil, errors.New("invalid or missing arguments")
		}
		return block.New("", args[1:]...), nil
	},
}

var Drop = shared.TaSignature{
	Name:       "drop",
	IsVariadic: false,
	Arguments: []block.Kind{
		block.BlockKind,
	},
	Returns:     block.BlockKind,
	Description: "Create a list containing all but the last item in the input list",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 1 {
			return nil, errors.New("invalid or missing arguments")
		}
		return block.New("", args[:argc-1]...), nil
	},
}

var Item = shared.TaSignature{
	Name:       "item",
	IsVariadic: false,
	Arguments: []block.Kind{
		block.BlockKind,
		block.DecimalKind,
	},
	Returns:     block.BlockKind,
	Description: "Returns a specific item from a list",
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		argc := len(args)
		if argc < 2 {
			return nil, errors.New("invalid or missing arguments")
		}

		if !args[1].IsDecimal() {
			return nil, errors.Errorf("`%s' is not an int", args[1].Text)
		}

		i, ok := args[1].Decimal.Int64()
		if !ok {
			return nil, errors.Errorf("`%s' is not an int", args[1].Text)
		}
		index := int(i)
		if index < 0 || index >= argc {
			return nil, errors.New("Out of bounds")
		}

		return args[0].Children[index], nil
	},
}
