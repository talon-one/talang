package time_test

import (
	"testing"

	"github.com/talon-one/talang/block"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestAfter(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`after 2006-01-02T19:04:05Z 2006-01-02T15:04:05Z`,
		nil,
		block.NewBool(true),
	}, helpers.Test{
		`after "2006-01-01T19:04:05Z" "2006-01-02T15:04:05Z"`,
		nil,
		block.NewBool(false),
	})
}

func TestBefore(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`before "2006-01-02T19:04:05Z" "2006-01-02T15:04:05Z"`,
		nil,
		block.NewBool(false),
	}, helpers.Test{
		`before "2006-01-01T19:04:05Z" "2006-01-02T15:04:05Z"`,
		nil,
		block.NewBool(true),
	})
}
