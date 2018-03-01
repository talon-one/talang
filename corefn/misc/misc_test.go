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
