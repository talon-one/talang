//go:generate go run ../generate_allop.go -pkg=mapping
package mapping

import (
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

var Map = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       "kv",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.BlockKind,
		},
		Returns:     block.MapKind,
		Description: "Create a map with any key value pairs passed as arguments.",
	},
	Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
		m := make(map[string]*block.Block)
		for i := 0; i < len(args); i++ {
			if args[i].IsBlock() && len(args[i].String) > 0 {
				if childCount := len(args[i].Children); childCount > 1 {
					children := make([]*block.Block, childCount)
					copy(children, args[i].Children)
					m[args[i].String] = block.NewList(children...)
				} else if childCount == 1 {
					m[args[i].String] = args[i].Children[0]
				}
			}
		}
		return block.NewMap(m), nil
	},
}
