//go:generate go run ../generate_allop.go -pkg=list
package list

import (
	"errors"

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
