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

func TestCatch(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`catch 22 (panic)`,
			nil,
			block.NewDecimalFromInt(22),
		},
		helpers.Test{
			`catch (+ 1 1) (panic)`,
			nil,
			block.NewDecimalFromInt(2),
		},
		helpers.Test{
			`catch "22" (. Profile)`,
			nil,
			block.NewString("22"),
		},
		helpers.Test{
			`catch 22 (. Profile)`,
			nil,
			block.NewDecimalFromInt(22),
		},
		helpers.Test{
			`catch 22 (+ 2 2)`,
			nil,
			block.NewDecimalFromInt(4),
		},
		helpers.Test{
			`catch 22 2`,
			nil,
			block.NewDecimalFromInt(2),
		},
		helpers.Test{
			`(catch (+ 2 (* 5 (- 3 4))) (+ 2 ( * 4 (panic))))`,
			nil,
			block.NewDecimalFromInt(-3),
		},
	)
}

func TestDo(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"do (list 1 2 3) Lst (. Lst)",
		nil,
		block.NewList(
			block.NewDecimalFromInt(1),
			block.NewDecimalFromInt(2),
			block.NewDecimalFromInt(3),
		),
	}, helpers.Test{
		"do 4 x (panic)",
		nil,
		helpers.Error{},
	})
}

func TestDoLegacy(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"do (list 1 2 3) ((Lst) (. Lst))",
		nil,
		block.NewList(
			block.NewDecimalFromInt(1),
			block.NewDecimalFromInt(2),
			block.NewDecimalFromInt(3),
		),
	}, helpers.Test{
		"do 4 ((x) (panic))",
		nil,
		helpers.Error{},
	}, helpers.Test{
		"do 4 (x (panic))",
		nil,
		helpers.Error{},
	}, helpers.Test{
		"do 4 (panic)",
		nil,
		helpers.Error{},
	})
}
