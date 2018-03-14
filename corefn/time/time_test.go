package time_test

import (
	"testing"

	"github.com/araddon/dateparse"

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
	// time, _ := time.Parse(time.RFC3339, "2018-01-02T19:04:05Z")
	time, _ := dateparse.ParseAny("2018-01-02T19:04:05Z")
	helpers.RunTests(t, helpers.Test{
		`parseTime 2018-01-02T19:04:05Z`,
		nil,
		block.NewTime(time),
	}, helpers.Test{
		`parseTime 2018-01-02T19:04:05Z`,
		nil,
		block.NewTime(time),
	})
}
