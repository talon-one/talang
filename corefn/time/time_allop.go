package time

import "github.com/talon-one/talang/interpreter"

func AllOperations() []interpreter.TaFunction {
	return []interpreter.TaFunction{
		After,
		Before,
		BetweenTimes,
		ParseTime,
		Date,
		Hour,
		Minute,
		Month,
		MonthDay,
		WeekDay,
		Year,
		FormatTime,
		MatchTime,
	}
}
