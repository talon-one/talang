//go:generate go run ../generate_allop.go -pkg=list
package list

import (
	"sort"
	"strings"

	"github.com/ericlagergren/decimal"
	"github.com/pkg/errors"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	if err := interpreter.RegisterCoreFunction(AllOperations()...); err != nil {
		panic(err)
	}
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
		Example: `
(list "Hello World" "Hello Universe")                            // returns a list with string items
(list 1 true Hello)                                              // returns a list with an int, bool and string
`,
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
		Example: `
(head (list "Hello World" "Hello Universe"))                     // returns "Hello World"
(head (list 1 true Hello))                                       // returns 1
`,
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
		Example: `
(tail (list "Hello World" "Hello Universe"))                     // returns a list containing "Hello Universe"
(tail (list 1 true Hello))                                       // returns a list containing true and Hello
`,
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
		Example: `
(drop (list "Hello World" "Hello Universe"))                     // returns a list containing "Hello World"
(drop (list 1 true Hello))                                       // returns a list containing 1 and true
`,
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
		Example: `
(item (list "Hello World" "Hello Universe") 0)                   // returns "Hello World"
(item (list 1 true Hello) 1)                                     // returns true
(item (list 1 true Hello) 3)                                     // fails
`,
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
		Example: `
(push (list "Hello World" "Hello Universe") "Hello Human")       // returns a list containing "Hello World", "Hello Universe" and "Hello Human"
(push (list 1 2) 3 4)                                            // returns a list containing 1, 2, 3 and 4
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewList(append(args[0].Children, args[1:]...)...), nil
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
		Example: `
(map  (list "World" "Universe") x (+ "Hello " (. x)))            // returns a list containing "Hello World" and "Hello Universe"
`,
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
			block.Copy(&result, blockToRun)
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
		Example: `
(sort  (list "World" "Universe"))                                // returns a list containing "Universe" and "World"
(sort  (list "World" "Universe") true)                           // returns a list containing "World" and "Universe"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		list := block.NewList()
		list.Children = make([]*block.Block, len(args[0].Children))
		copy(list.Children, args[0].Children)

		if len(args) > 1 && args[1].Bool {
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
		Example: `
(min  (list 3 4 1 3 7 1 17 15 2))                                // returns 1
`,
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
		Example: `
(max  (list 3 4 1 3 7 1 17 15 2))                                // returns 17
`,
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

var Append = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:        "append",
		IsVariadic:  Push.IsVariadic,
		Arguments:   Push.Arguments,
		Returns:     Push.Returns,
		Description: Push.Description,
		Example:     Push.Example,
	},
	Func: Push.Func,
}

var Count = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "count",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.DecimalKind,
		Description: "Return the number of items in the input list",
		Example: `
(count (list 1 2 3 4))											 // returns "4"
(count (list 1))												 // returns "1"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		count := int64(len(args[0].Children))
		return block.NewDecimalFromInt(count), nil
	},
}

var Reverse = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "reverse",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.ListKind,
		Description: "Reverses the order of items in a given list",
		Example: `
(reverse (list 1 2 3 4))										 // returns "(4 3 2 1)"
(reverse (list 1))												 // returns "(1)"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		list := block.NewList()
		list.Children = make([]*block.Block, len(args[0].Children))
		childrenCount := len(args[0].Children) - 1
		for i := childrenCount; i >= 0; i-- {
			list.Children[childrenCount-i] = args[0].Children[i]
		}
		return list, nil
	},
}

var Join = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "join",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
			block.StringKind,
		},
		Returns:     block.StringKind,
		Description: "Create a string by joining together a list of strings with `glue`",
		Example: `
(join (list hello world) "-")									 // returns "hello-world"
(join (list hello world) ",")									 // returns "hello,world"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		if args[0].Children[0].Kind != block.StringKind {
			return nil, errors.New("List must be of string type")
		}

		var final = make([]string, len(args[0].Children))
		for i := 0; i < len(args[0].Children); i++ {
			final[i] = args[0].Children[i].String
		}
		ret := strings.Join(final, args[1].String)
		return block.NewString(ret), nil
	},
}

var IsEmpty = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "isEmpty",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.ListKind,
		},
		Returns:     block.BoolKind,
		Description: "Check if a list is empty",
		Example: `
isEmpty (list hello world)				                         // returns "false"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		if args[0].IsEmpty() {
			return block.NewBool(true), nil
		}
		return block.NewBool(false), nil
	},
}

var Split = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "split",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.StringKind,
			block.StringKind,
		},
		Returns:     block.ListKind,
		Description: "Create a list of strings by splitting the given string at each occurence of `sep`",
		Example: `
(split "1,2,3,a" ",")				                             // returns "[1 2 3 a]"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		list := block.NewList()
		srcs := strings.Split(args[0].String, args[1].String)
		list.Children = make([]*block.Block, len(srcs))
		for i := 0; i < len(list.Children); i++ {
			list.Children[i] = block.NewString(srcs[i])
		}
		return list, nil
	},
}
