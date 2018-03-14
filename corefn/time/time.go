//go:generate go run ../generate_allop.go -pkg=time
package time

import (
	"github.com/araddon/dateparse"
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

var BetweenTimes = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "betweenTimes",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind, // timestamp
			block.TimeKind, // minTime
			block.TimeKind, // maxTime
		},
		Returns:     block.BoolKind,
		Description: "Evaluates whether a timestamp is between minTime and maxTime",
		Example: `
(betweenTimes "2006-01-02T19:04:05Z" "2006-01-01T15:04:05Z" "2006-01-03T19:04:05Z")                                // returns "false"
(betweenTimes "2006-01-01T19:04:05Z" "2006-01-02T15:04:05Z" "2006-01-03T19:04:05Z")                                // returns "true"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		a := args[0].Time.After(args[1].Time)
		b := args[0].Time.Before(args[2].Time)
		return block.NewBool(a && b), nil
	},
}

var ParseTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "parseTime",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.TimeKind,
		Description: "Evaluates whether a timestamp is between minTime and maxTime",
		Example:     ``,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		date, err := dateparse.ParseAny(args[0].String)
		if err != nil {
			return nil, err
		}
		return block.NewTime(date), nil
	},
}
