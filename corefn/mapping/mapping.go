//go:generate go run ../generate_allop.go -pkg=mapping
package mapping

import (
	"errors"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	interpreter.RegisterCoreFunction(AllOperations()...)
}

var KV = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "kv",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.BlockKind,
		},
		Returns:     block.MapKind,
		Description: "Create a map with any key value pairs passed as arguments.",
		Example: `
(kv (Key1 "Hello World") (Key2 true) (Key3 123))               // returns a Map with the keys key1, key2, key3
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		m := make(map[string]*block.Block)
		for i := 0; i < len(args); i++ {
			if args[i].IsBlock() && len(args[i].String) > 0 {
				if len(args[i].Children) > 1 {
					return nil, errors.New("Unable to add multiple values to one key")
				}
				value := args[i].Children[0]
				if value.IsBlock() {
					if err := interp.Evaluate(value); err != nil {
						return nil, err
					}
				}
				m[args[i].String] = value
			}
		}
		return block.NewMap(m), nil
	},
}
