//go:generate go run ../generate_allop.go -pkg=time
package time

import (
	"fmt"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/token"
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
		Arguments: []token.Kind{
			token.Time,
			token.Time,
		},
		Returns:     token.Bool,
		Description: "Checks whether time A is after B",
		Example: `
(after 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z)               ; returns "true"
(after 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z)               ; returns "false"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		time := args[0].Time.After(args[1].Time)
		return token.NewBool(time), nil
	},
}

var Before = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "before",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
			token.Time,
		},
		Returns:     token.Bool,
		Description: "Checks whether time A is before B",
		Example: `
(before 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z)              ; returns "false"
(before 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z)              ; returns "true"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		time := args[0].Time.Before(args[1].Time)
		return token.NewBool(time), nil
	},
}

var BetweenTimes = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "betweenTimes",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time, // timestamp
			token.Time, // minTime
			token.Time, // maxTime
		},
		Returns:     token.Bool,
		Description: "Evaluates whether a timestamp is between minTime and maxTime",
		Example: `
(betweenTimes 2006-01-02T19:04:05Z 2006-01-01T15:04:05Z 2006-01-03T19:04:05Z)                                ; returns "false"
(betweenTimes 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z 2006-01-03T19:04:05Z)                                ; returns "true"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		a := args[0].Time.After(args[1].Time)
		b := args[0].Time.Before(args[2].Time)
		return token.NewBool(a && b), nil
	},
}

var ParseTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "parseTime",
		IsVariadic: true,
		Arguments: []token.Kind{
			token.String, // time string
			token.String, // signature
		},
		Returns:     token.Time,
		Description: "Evaluates whether a timestamp is between minTime and maxTime",
		Example: `
(parseTime "2018-01-02T19:04:05Z")                               ; returns "2018-01-02 19:04:05 +0000 UTC"
(parseTime "20:04:05Z" "HH:mm:ss")                               ; returns "2018-01-02 20:04:05 +0000 UTC"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
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
		return token.NewTime(date), nil
	},
}

var Hour = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "hour",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.String,
		Description: "Extract the hour (00-23) from a time",
		Example: `
(hour 2018-01-14T19:04:05Z)                                      ; returns "19"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		hour := strconv.Itoa(args[0].Time.Hour())
		return token.NewString(hour), nil
	},
}

var Minute = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "minute",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.String,
		Description: "Extract the hour (00-23) from a time",
		Example: `
(minute 2018-01-14T19:04:05Z)                                    ; returns "04"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		minute := strconv.Itoa(args[0].Time.Minute())
		return token.NewString(minute), nil
	},
}

var Date = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "date",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.String,
		Description: "Extract the date in YYYY-MM-DD format from a time.",
		Example: `
(betweenTimes 2006-01-02T19:04:05Z 2006-01-01T15:04:05Z 2006-01-03T19:04:05Z)                                ; returns "false"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		return token.NewString(jodaTime.Format("yyyy-MM-dd", args[0].Time)), nil
	},
}

var Month = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "month",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.String,
		Description: "Extract the month (1-12) from a time",
		Example: `
(month 2018-01-02T19:04:05Z)                                     ; returns "1"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		month := strconv.Itoa(int(args[0].Time.Month()))
		return token.NewString(month), nil
	},
}

var MonthDay = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "monthDay",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.String,
		Description: "Extract the day (1-31) from a time",
		Example: `
(monthDay 2018-01-14T19:04:05Z)                                  ; returns "14"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		monthDay := strconv.Itoa(int(args[0].Time.Day()))
		return token.NewString(monthDay), nil
	},
}

// Disclaimer: weekDay has lowercased 'D' due to old names
var WeekDay = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "weekday",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.String,
		Description: "Extract the week day (0-6) from a time",
		Example: `
(weekDay 2018-01-14T19:04:05Z)                                   ; returns "3"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		return token.NewString(jodaTime.Format("e", args[0].Time)), nil
	},
}

var Year = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "year",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.String,
		Description: "Extract the year from a time",
		Example: `
(year 2018-01-02T19:04:05Z)                                      ; returns "2018"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		year := strconv.Itoa(int(args[0].Time.Year()))
		return token.NewString(year), nil
	},
}

var FormatTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "formatTime",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.String,
		Description: "Create an RFC3339 timestamp, the inverse of parseTime",
		Example: `
(formatTime 2018-01-02T19:04:05Z)                                ; returns "2018"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		return token.NewString(args[0].Time.Format(time.RFC3339)), nil
	},
}

var MatchTime = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "matchTime",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,   // timestamp1
			token.Time,   // timestamp1
			token.String, // layout
		},
		Returns:     token.Bool,
		Description: "Checks if two times match for a given layout",
		Example: `
matchTime 2018-03-11T00:04:05Z 2018-03-11T00:04:05Z YYYY-MM-DD   ; returns "true"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		layout := args[2].String
		return token.NewBool(jodaTime.Format(layout, args[0].Time) == jodaTime.Format(layout, args[1].Time)), nil
	},
}

// TODO: test coverage. Days returns a float64 to be consistent, this need to be dealt with.
var Days = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "days",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,
		},
		Returns:     token.Decimal,
		Description: "Extract days from now from time",
		Example: `
(days 2018-03-18T00:04:05Z)                                      ; returns "3.423892107645601701193527333089150488376617431640625" results vary as the function is relative to the current date.
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		now := args[0].Time.Sub(time.Now())
		return token.NewDecimalFromFloat((now.Hours() / 24)), nil
	},
}

var AddDuration = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "addDuration",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,    // since
			token.Decimal, // amount
			token.String,  // units
		},
		Returns:     token.Time,
		Description: "Extract days from now from time",
		Example: `
(addDuration 2018-03-18T00:04:05Z 3 minutes)                     ; returns "2018-03-18T00:07:05Z"
(addDuration 2018-03-18T00:04:05Z 2 hours)                       ; returns "2018-03-18T02:04:05Z"
(addDuration 2018-03-18T00:04:05Z 18 days)                       ; returns "2018-04-05T00:04:05Z"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		duration, err := makeDuration(args[1], args[2].String)
		fmt.Println(duration)
		if err != nil {
			return nil, err
		}
		return token.NewTime(args[0].Time.Add(duration)), nil
	},
}

var SubDuration = interpreter.TaFunction{
	CommonSignature: interpreter.CommonSignature{
		Name:       "subDuration",
		IsVariadic: false,
		Arguments: []token.Kind{
			token.Time,    // since
			token.Decimal, // amount
			token.String,  // units
		},
		Returns:     token.Time,
		Description: "Extract days from now from time",
		Example: `
(subDuration 2018-03-18T00:04:05Z 12 minutes)                    ; returns "2018-03-17T23:52:05Z"
(subDuration 2018-03-18T00:04:05Z 17 hours)                      ; returns "2018-03-17T07:04:05Z"
(subDuration 2018-03-18T00:04:05Z 22 days)                       ; returns "2018-02-24T00:04:05Z"
`,
	},
	Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
		duration, err := makeDuration(args[1], args[2].String)
		if err != nil {
			return nil, err
		}
		return token.NewTime(args[0].Time.Add(-1 * duration)), nil
	},
}

func makeDuration(n *token.TaToken, unit string) (time.Duration, error) {
	var multiplier int64
	switch unit {
	case "days":
		multiplier = int64(time.Hour) * 24
	case "hours":
		multiplier = int64(time.Hour)
	case "minutes":
		multiplier = int64(time.Minute)
	default:
		return 0, fmt.Errorf("invalid duration unit %q", unit)
	}
	trg, _ := n.Decimal.Int64()
	result := multiplier * trg
	return time.Duration(result), nil
}
