//go:generate go run ../generate_allop.go -pkg=time
package time

import (
	"time"

	"github.com/araddon/dateparse"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
	"github.com/vjeantet/jodaTime"
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
(after 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z)                                // returns "true"
(after 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z)                                // returns "false"
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
(before 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z)                                // returns "false"
(before 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z)                                // returns "true"
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
(betweenTimes 2006-01-02T19:04:05Z 2006-01-01T15:04:05Z 2006-01-03T19:04:05Z)                                // returns "false"
(betweenTimes 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z 2006-01-03T19:04:05Z)                                // returns "true"
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
		IsVariadic: true,
		Arguments: []block.Kind{
			block.StringKind, // time string
			block.StringKind, // signature
		},
		Returns:     block.TimeKind,
		Description: "Evaluates whether a timestamp is between minTime and maxTime",
		Example: `
(parseTime "2018-01-02T19:04:05Z")                              // returns "2018-01-02 19:04:05 +0000 UTC"
(parseTime "20:04:05Z" "HH:mm:ss")                              // returns "2018-01-02 20:04:05 +0000 UTC"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		var date time.Time
		var err error
		if len(args) > 1 {
			date, err = jodaTime.Parse(args[1].String, args[0].String)
			if err != nil {
				return nil, err
			}
		} else {
			date, err = dateparse.ParseAny(args[0].String)
			if err != nil {
				return nil, err
			}
		}
		return block.NewTime(date), nil
	},
}

var Hour = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "hour",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Extract the hour (00-23) from a time",
		Example: `
(hour 2018-01-14T19:04:05Z)                                // returns "19"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(jodaTime.Format("HH", args[0].Time)), nil
	},
}

var Date = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "date",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Extract the date in YYYY-MM-DD format from a time.",
		Example: `
(betweenTimes 2006-01-02T19:04:05Z 2006-01-01T15:04:05Z 2006-01-03T19:04:05Z)                                // returns "false"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(jodaTime.Format("yyyy-MM-dd", args[0].Time)), nil
	},
}

var Month = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "month",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Extract the month (1-11) from a time",
		Example: `
(month 2018-01-02T19:04:05Z)                                // returns "1"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(jodaTime.Format("M", args[0].Time)), nil
	},
}

var MonthDay = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "monthDay",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Extract the day (1-31) from a time",
		Example: `
(monthDay 2018-01-14T19:04:05Z)                                // returns "14"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(jodaTime.Format("d", args[0].Time)), nil
	},
}

// Disclaimer: weekDay has lowercased 'D' due to old names
var WeekDay = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "weekday",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Extract the week day (0-6) from a time",
		Example: `
(weekDay 2018-01-14T19:04:05Z)                                // returns "3"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(jodaTime.Format("e", args[0].Time)), nil
	},
}

var Year = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "year",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Extract the year from a time",
		Example: `
(year 2018-01-02T19:04:05Z)                                // returns "2018"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(jodaTime.Format("yyyy", args[0].Time)), nil
	},
}

var FormatTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "formatTime",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.TimeKind,
		},
		Returns:     block.StringKind,
		Description: "Create an RFC3339 timestamp, the inverse of parseTime",
		Example: `
(formatTime 2018-01-02T19:04:05Z)                                // returns "2018"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
		return block.NewString(args[0].Time.Format(time.RFC3339)), nil
	},
}
