package math_test

import (
	"testing"

	"github.com/talon-one/talang/lexer"
	"github.com/talon-one/talang/token"

	helpers "github.com/talon-one/talang/testhelpers"
)

func TestAdd(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"+ 1",
			nil,
			lexer.MustLex("+ 1"),
		},
		helpers.Test{
			"+ 1 2",
			nil,
			token.NewDecimalFromInt(3),
		},
		helpers.Test{
			"+ 1.2 3.4",
			nil,
			token.NewDecimalFromString("4.6"),
		},
		helpers.Test{
			"+ 1.2 3.4 -5.6",
			nil,
			token.NewDecimalFromInt(-1),
		},
	)
}

func TestSub(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"- 1",
			nil,
			lexer.MustLex("- 1"),
		},
		helpers.Test{
			"- 1 2",
			nil,
			token.NewDecimalFromInt(-1),
		},
		helpers.Test{
			"- 1.2 3.4",
			nil,
			token.NewDecimalFromString("-2.2"),
		},
		helpers.Test{
			"- 1.2 3.4 -5.6",
			nil,
			token.NewDecimalFromString("3.4"),
		},
	)
}

func TestMul(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"* 1",
			nil,
			lexer.MustLex("* 1"),
		},
		helpers.Test{
			"* 1 2",
			nil,
			token.NewDecimalFromInt(2),
		},
		helpers.Test{
			"* 1.2 3.4",
			nil,
			token.NewDecimalFromString("4.08"),
		},
		helpers.Test{
			"* 1.2 3.4 -5.6",
			nil,
			token.NewDecimalFromString("-22.848"),
		},
	)
}

func TestDiv(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"/ 1",
			nil,
			lexer.MustLex("/ 1"),
		},
		helpers.Test{
			"/ 1 2",
			nil,
			token.NewDecimalFromString("0.5"),
		},
		helpers.Test{
			"/ 1.2 3.4",
			nil,
			token.NewDecimalFromString("0.3529411764705882"),
		},
		helpers.Test{
			"/ 1.2 3.4 -5.6",
			nil,
			token.NewDecimalFromString("-0.06302521008403361"),
		},
	)
}

func TestMod(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"mod 1",
			nil,
			lexer.MustLex("mod 1"),
		},
		helpers.Test{
			"mod 2 1",
			nil,
			token.NewDecimalFromInt(0),
		},
		helpers.Test{
			"mod 3 2",
			nil,
			token.NewDecimalFromInt(1),
		},
		helpers.Test{
			"mod 4 3 2",
			nil,
			token.NewDecimalFromInt(1),
		},
	)
}

func TestFloor(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"floor 2",
			nil,
			token.NewDecimalFromInt(2),
		},
		helpers.Test{
			"floor 2.4",
			nil,
			token.NewDecimalFromInt(2),
		},
		helpers.Test{
			"floor 2.5",
			nil,
			token.NewDecimalFromInt(2),
		},
		helpers.Test{
			"floor 2.9",
			nil,
			token.NewDecimalFromInt(2),
		},
		helpers.Test{
			"floor -2.7",
			nil,
			token.NewDecimalFromInt(-3),
		},
		helpers.Test{
			"floor -2",
			nil,
			token.NewDecimalFromInt(-2),
		},
	)
}

func TestCeil(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"ceil 2",
			nil,
			token.NewDecimalFromInt(2),
		},
		helpers.Test{
			"ceil 2.4",
			nil,
			token.NewDecimalFromInt(3),
		},
		helpers.Test{
			"ceil 2.5",
			nil,
			token.NewDecimalFromInt(3),
		},
		helpers.Test{
			"ceil 2.9",
			nil,
			token.NewDecimalFromInt(3),
		},
		helpers.Test{
			"ceil -2.7",
			nil,
			token.NewDecimalFromInt(-2),
		},
		helpers.Test{
			"ceil -2",
			nil,
			token.NewDecimalFromInt(-2),
		},
	)
}
