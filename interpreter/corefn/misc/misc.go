package misc

import (
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var Misc = shared.TaSignature{
	Name: "misc3",
	Arguments: []block.Kind{
		block.BlockKind,
	},
	Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
		return args[0], nil
	},
}
