//go:generate go run ../generate_allop.go -pkg=list
package list

import (
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
