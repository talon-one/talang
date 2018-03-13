package misc_test

import (
	"testing"

	"github.com/talon-one/talang/block"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestNoop(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"noop",
		nil,
		block.NewNull(),
	})
}

func TestToString(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"toString 1",
		nil,
		block.NewString("1"),
	})
}

func TestNot(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"not false",
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			"not true",
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			"not (not false)",
			nil,
			block.NewBool(false),
		},
	)
}
