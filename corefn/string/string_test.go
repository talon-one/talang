package string_test

import (
	"testing"

	"github.com/talon-one/talang/token"
	"github.com/talon-one/talang/lexer"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestAdd(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`+ "Hello World" " and " Universe`,
		nil,
		token.NewString("Hello World and Universe"),
	})
}

func TestContains(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`contains "Hello"`,
			nil,
			lexer.MustLex(`contains "Hello"`),
		},
		helpers.Test{
			`contains "Hello World" Universe`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`contains "Hello World" Hello Universe`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`contains "Hello World" Hello World`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestNotContains(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`notContains "Hello"`,
			nil,
			lexer.MustLex(`notContains "Hello"`),
		},
		helpers.Test{
			`notContains "Hello World" Universe`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`notContains "Hello World" Hello Universe`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`notContains "Hello World" Hello World`,
			nil,
			token.NewBool(false),
		},
	)
}

func TestStartsWith(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`startsWith "Hello"`,
			nil,
			lexer.MustLex(`startsWith "Hello"`),
		},
		helpers.Test{
			`startsWith "Hello World" Bye`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`startsWith "Hello World" Hello Bye`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`startsWith "Hello World" Hello Hell`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestEndsWith(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`endsWith "Hello"`,
			nil,
			lexer.MustLex(`endsWith "Hello"`),
		},
		helpers.Test{
			`endsWith "Hello World" Universe`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`endsWith "Hello World" World Universe`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`endsWith "Hello World" World ld`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestRegexp(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`~ "Hello"`,
			nil,
			lexer.MustLex(`~ "Hello"`),
		},
		helpers.Test{
			`~ [a foo`,
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			`~ ^foo foobar`,
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			`~ ^foo$ foobar`,
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			`~ "^Hello\s\w+" "Hello World" "Hello Universe"`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestLastName(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`lastName "Hello Mr Mock"`,
			nil,
			token.NewString("Mock"),
		}, helpers.Test{
			`lastName "Bond"`,
			nil,
			token.NewString("Bond"),
		},
	)
}

func TestFirstName(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`firstName "Hello Mr Mock"`,
			nil,
			token.NewString("Hello"),
		}, helpers.Test{
			`firstName "Bond"`,
			nil,
			token.NewString("Bond"),
		},
	)
}
