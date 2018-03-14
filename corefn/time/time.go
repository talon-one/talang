//go:generate go run ../generate_allop.go -pkg=time
package time

import (
	_ "github.com/araddon/dateparse"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func init() {
	if err := interpreter.RegisterCoreFunction(AllOperations()...); err != nil {
		panic(err)
	}
}

var After = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "after",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Checks whether time A is after B",
		Example: `
(after "2006-01-02T19:04:05Z" "2006-01-02T15:04:05Z")                                // returns "true"
(after "2006-01-01T19:04:05Z" "2006-01-02T15:04:05Z")                                // returns "false"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		time := args[0].Time.After(args[1].Time)
		return block.NewBool(time), nil
	},
}

var Before = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "before",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
			block.TimeKind,
		},
		Returns:     block.BoolKind,
		Description: "Checks whether time A is before B",
		Example: `
(before "2006-01-02T19:04:05Z" "2006-01-02T15:04:05Z")                                // returns "false"
(before "2006-01-01T19:04:05Z" "2006-01-02T15:04:05Z")                                // returns "true"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		time := args[0].Time.Before(args[1].Time)
		return block.NewBool(time), nil
	},
}
