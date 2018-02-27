package string_test

import (
	"testing"

	"github.com/talon-one/talang/block"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestAdd(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`+ "Hello World" " and " Universe`,
		nil,
		block.NewString("Hello World and Universe"),
	})
}

func TestContains(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`contains "Hello"`,
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			`contains "Hello World" Universe`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`contains "Hello World" Hello Universe`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`contains "Hello World" Hello World`,
			nil,
			block.NewBool(true),
		},
	)
}

func TestNotContains(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`notContains "Hello"`,
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			`notContains "Hello World" Universe`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`notContains "Hello World" Hello Universe`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`notContains "Hello World" Hello World`,
			nil,
			block.NewBool(false),
		},
	)
}

func TestStartsWith(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`startsWith "Hello"`,
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			`startsWith "Hello World" Bye`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`startsWith "Hello World" Hello Bye`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`startsWith "Hello World" Hello Hell`,
			nil,
			block.NewBool(true),
		},
	)
}

func TestEndsWith(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`endsWith "Hello"`,
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			`endsWith "Hello World" Universe`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`endsWith "Hello World" World Universe`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`endsWith "Hello World" World ld`,
			nil,
			block.NewBool(true),
		},
	)
}

func TestRegexp(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`~ "Hello"`,
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			`~ [a foo`,
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			`~ ^foo foobar`,
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			`~ ^foo$ foobar`,
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			`~ "^Hello\s\w+" "Hello World" "Hello Universe"`,
			nil,
			block.NewBool(true),
		},
	)
}
