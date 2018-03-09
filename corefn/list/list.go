//go:generate go run ../generate_allop.go -pkg=list
package list

import (
	"sort"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	interpreter.RegisterCoreFunction(AllOperations()...)
}

var List = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "list",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AtomKind,
			block.AtomKind,
		},
		Returns:     block.ListKind,
		Description: "Create a list out of the children",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewList(args...), nil
	},
}

var Head = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "head",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.AnyKind,
		Description: "Returns the first item in the list",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		if len(args[0].Children) > 0 {
			return args[0].Children[0], nil
		}
		return block.NewNull(), nil
	},
}

var Tail = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "tail",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.ListKind,
		Description: "Returns list without the first item",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		if len(args[0].Children) <= 0 {
			return block.NewList(), nil
		}
		return block.NewList(args[0].Children[1:]...), nil
	},
}

var Drop = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "drop",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.ListKind,
		Description: "Create a list containing all but the last item in the input list",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		if l := len(args[0].Children); l > 0 {
			return block.NewList(args[0].Children[:l-1]...), nil
		}
		return block.NewList(), nil
	},
}

var Item = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "item",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
			block.DecimalKind,
		},
		Returns:     block.AnyKind,
		Description: "Returns a specific item from a list",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
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

var Push = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "push",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.ListKind,
			block.AtomKind | block.CollectionKind,
			block.AtomKind | block.CollectionKind,
		},
		Returns:     block.ListKind,
		Description: "Adds an item to the list and returns the list",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var list block.Block
		list.Update(args[0])
		list.Children = append(list.Children, args[1:]...)
		return &list, nil
	},
}

var Map = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name: "map",
		Arguments: []block.Kind{
			block.ListKind,
			block.StringKind,
			block.BlockKind,
		},
		Returns:     block.ListKind,
		Description: "Create a new list by evaluating the given block for each item in the input list",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		list := args[0]
		bindingName := args[1].String
		blockToRun := args[2]

		size := len(list.Children)
		values := make([]*block.Block, 0, size)
		scope := interp.NewScope()

		for i := 0; i < size; i++ {
			scope.Set(bindingName, list.Children[i])

			var result block.Block
			result.Update(blockToRun)
			if err := scope.Evaluate(&result); err != nil {
				return nil, err
			}
			values = append(values, &result)
		}
		return block.NewList(values...), nil
	},
}

var Sort = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "sort",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.ListKind,
			block.BoolKind,
		},
		Returns:     block.ListKind,
		Description: "Sort a list ascending, set the second argument to true for descending order",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		list := block.NewList()
		list.Children = make([]*block.Block, len(args[0].Children), len(args[0].Children))
		copy(list.Children, args[0].Children)

		if len(args) > 1 && args[1].Bool == true {
			a := block.BlockArguments(list.Children)
			sort.Sort(sort.Reverse(&a))
		} else {
			sort.Sort(block.BlockArguments(list.Children))
		}

		return list, nil
	},
}

var Min = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "min",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.DecimalKind,
		Description: "Find the lowest number in the list",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var d *decimal.Big
		for _, item := range args[0].Children {
			if item.IsDecimal() {
				if d == nil {
					d = item.Decimal
				} else if item.Decimal.Cmp(d) < 0 {
					d = item.Decimal
				}
			}
		}
		if d == nil {
			return nil, errors.New("No decimal present in list")
		}
		return block.NewDecimal(d), nil
	},
}

var Max = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "max",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.DecimalKind,
		Description: "Find the largest number in the list",
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var d *decimal.Big
		for _, item := range args[0].Children {
			if item.IsDecimal() {
				if d == nil {
					d = item.Decimal
				} else if item.Decimal.Cmp(d) > 0 {
					d = item.Decimal
				}
			}
		}
		if d == nil {
			return nil, errors.New("No decimal present in list")
		}
		return block.NewDecimal(d), nil
	},
}
