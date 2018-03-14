package time_test

import (
	"testing"
	"time"

	"github.com/araddon/dateparse"
	"github.com/vjeantet/jodaTime"

	"github.com/talon-one/talang/block"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestAfter(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`after 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z`,
		nil,
		block.NewBool(true),
	}, helpers.Test{
		`after 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z`,
		nil,
		block.NewBool(false),
	})
}

func TestBefore(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`before 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z`,
		nil,
		block.NewBool(false),
	}, helpers.Test{
		`before 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z`,
		nil,
		block.NewBool(true),
	})
}

func TestBetweenTimes(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`betweenTimes 2006-01-02T19:04:05Z 2006-01-01T15:04:05Z 2006-01-03T19:04:05Z`,
		nil,
		block.NewBool(true),
	}, helpers.Test{
		`betweenTimes 2006-01-01T19:04:05Z 2006-01-02T15:04:05Z 2006-01-03T19:04:05Z`,
		nil,
		block.NewBool(false),
	})
}

func TestParseTime(t *testing.T) {
	_time, _ := dateparse.ParseAny("2018-01-02T19:04:05Z")
	helpers.RunTests(t, helpers.Test{
		`parseTime "2018-01-02T19:04:05Z"`,
		nil,
		block.NewTime(_time),
	}, helpers.Test{
		`parseTime "2018-01-02T19:04:05Z"`,
		nil,
		block.NewTime(_time),
	}, helpers.Test{
		`parseTime "-42"`,
		nil,
		helpers.Error{},
	}, helpers.Test{
		`parseTime 10:30:31 HH:mm:ss`,
		nil,
		block.NewTime(mustParseJodaTime("HH:mm:ss", "10:30:31")),
	}, helpers.Test{
		`parseTime 10:30:31 OO:TT:{{`,
		nil,
		helpers.Error{},
	})
}

func mustParseJodaTime(layout string, date string) time.Time {
	time, err := jodaTime.Parse(layout, date)
	if err != nil {
		panic(err)
	}
	return time
}

func TestDate(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`date 2018-01-02T19:04:05Z`,
		nil,
		block.NewString("2018-01-02"),
	})
}

func TestMonth(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`month 2018-01-02T19:04:05Z`,
		nil,
		block.NewString("1"),
	})
}

func TestYear(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`year 2018-01-02T19:04:05Z`,
		nil,
		block.NewString("2018"),
	})
}

func TestMonthDay(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`monthDay 2018-01-14T19:04:05Z`,
		nil,
		block.NewString("14"),
	})
}
func TestWeekDay(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`weekDay 2018-03-11T19:04:05Z`,
		nil,
		block.NewString("0"),
	})
}

func TestFormatTime(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`formatTime 2018-03-11T19:04:05Z`,
		nil,
		block.NewString("2018-03-11T19:04:05Z"),
	})
}

func TestHour(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`hour 2018-03-11T00:04:05Z`,
		nil,
		block.NewString("0"),
	}, helpers.Test{
		`hour 2018-03-11T23:04:05Z`,
		nil,
		block.NewString("23"),
	}, helpers.Test{
		`hour 2018-03-11T04:04:05Z`,
		nil,
		block.NewString("4"),
	})
}

func TestMinute(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`minute 2018-03-11T00:04:05Z`,
		nil,
		block.NewString("4"),
	}, helpers.Test{
		`minute 2018-03-11T00:52:05Z`,
		nil,
		block.NewString("52"),
	})
}

func TestMatchTime(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`matchTime 2018-03-11T00:04:05Z 2018-03-11T00:04:05Z YYYY-MM-DD`,
		nil,
		block.NewBool(true),
	}, helpers.Test{
		`matchTime 2018-04-11T00:04:05Z 2018-03-11T00:04:05Z YYYY-MM-DD`,
		nil,
		block.NewBool(false),
	})
}

func TestAddDuration(t *testing.T) {
	_time1, _ := time.Parse(time.RFC3339, "2018-03-11T00:05:05Z")
	_time2, _ := time.Parse(time.RFC3339, "2018-03-12T02:04:05Z")
	_time3, _ := time.Parse(time.RFC3339, "2018-03-19T02:04:05Z")
	helpers.RunTests(t, helpers.Test{
		`addDuration 2018-03-11T00:04:05Z 1 "minutes"`,
		nil,
		block.NewTime(_time1),
	}, helpers.Test{
		`addDuration 2018-03-12T00:04:05Z 2 "hours"`,
		nil,
		block.NewTime(_time2),
	}, helpers.Test{
		`addDuration 2018-03-14T02:04:05Z 5 "days"`,
		nil,
		block.NewTime(_time3),
	}, helpers.Test{
		`addDuration 2018-03-14T02:04:05Z 5 "eons"`,
		nil,
		helpers.Error{},
	})
}
