//go:generate go run ../generate_allop.go -pkg=mapping
package mapping

import (
	"errors"

	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/token"
)

func init() {
	if err := interpreter.RegisterCoreFunction(AllOperations()...); err != nil {
		panic(err)
	}
}

var KV = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "kv",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.Token,
		},
		Returns:     token.Map,
		Description: "Create a map with any key value pairs passed as arguments.",
		Example: `
(kv (Key1 "Hello World") (Key2 true) (Key3 123))                 ; returns a Map with the keys key1, key2, key3
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		m := make(map[string]*token.TaToken)
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
		return token.NewMap(m), nil
	},
}
