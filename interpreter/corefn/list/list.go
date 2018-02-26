//go:generate go run ../generate_allop.go -pkg=list
package list

import (
	"github.com/pkg/errors"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var List = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "list",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AtomKind,
		},
		Returns:     block.ListKind,
		Description: "Create a list out of the children",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewList(args...), nil
	},
}

var Head = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "head",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.AnyKind,
		Description: "Returns the first item in the list",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		if len(args[0].Children) > 0 {
			return args[0].Children[0], nil
		}
		return block.NewNull(), nil
	},
}

var Tail = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "tail",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.ListKind,
		Description: "Returns list without the first item",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		if len(args[0].Children) <= 0 {
			return block.NewList(), nil
		}
		return block.NewList(args[0].Children[1:]...), nil
	},
}

var Drop = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "drop",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.ListKind,
		Description: "Create a list containing all but the last item in the input list",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		if l := len(args[0].Children); l > 0 {
			return block.NewList(args[0].Children[:l-1]...), nil
		}
		return block.NewList(), nil
	},
}

var Item = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "item",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
			block.DecimalKind,
		},
		Returns:     block.AnyKind,
		Description: "Returns a specific item from a list",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		i, ok := args[1].Decimal.Int64()
		if !ok {
			return nil, errors.Errorf("`%s' is not an int", args[1].String)
		}
		index := int(i)
		l := len(args[0].Children)
		if index < 0 || index >= l {
			return nil, errors.New("Out of bounds")
		}

		return args[0].Children[index], nil
	},
}
