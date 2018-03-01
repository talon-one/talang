package math_test

import (
	"testing"

	"github.com/talon-one/talang/lexer"

	"github.com/talon-one/talang/block"

	"github.com/ericlagergren/decimal"

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
			block.NewDecimal(decimal.New(3, 0)),
		},
		helpers.Test{
			"+ 1.2 3.4",
			nil,
			block.NewDecimal(decimal.New(46, 1)),
		},
		helpers.Test{
			"+ 1.2 3.4 -5.6",
			nil,
			block.NewDecimal(decimal.New(-1, 0)),
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
			block.NewDecimal(decimal.New(-1, 0)),
		},
		helpers.Test{
			"- 1.2 3.4",
			nil,
			block.NewDecimal(decimal.New(-22, 1)),
		},
		helpers.Test{
			"- 1.2 3.4 -5.6",
			nil,
			block.NewDecimal(decimal.New(34, 1)),
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
			block.NewDecimal(decimal.New(2, 0)),
		},
		helpers.Test{
			"* 1.2 3.4",
			nil,
			block.NewDecimal(decimal.New(408, 2)),
		},
		helpers.Test{
			"* 1.2 3.4 -5.6",
			nil,
			block.NewDecimal(decimal.New(-22848, 3)),
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
			block.NewDecimal(decimal.New(5, 1)),
		},
		helpers.Test{
			"/ 1.2 3.4",
			nil,
			block.NewDecimal(decimal.New(3529411764705882, 16)),
		},
		helpers.Test{
			"/ 1.2 3.4 -5.6",
			nil,
			block.NewDecimal(decimal.New(-6302521008403361, 17)),
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
			block.NewDecimal(decimal.New(0, 0)),
		},
		helpers.Test{
			"mod 3 2",
			nil,
			block.NewDecimal(decimal.New(1, 0)),
		},
		helpers.Test{
			"mod 4 3 2",
			nil,
			block.NewDecimal(decimal.New(1, 0)),
		},
	)
}

func TestFloor(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"floor 2",
			nil,
			block.NewDecimal(decimal.New(2, 0)),
		},
		helpers.Test{
			"floor 2.4",
			nil,
			block.NewDecimal(decimal.New(2, 0)),
		},
		helpers.Test{
			"floor 2.5",
			nil,
			block.NewDecimal(decimal.New(2, 0)),
		},
		helpers.Test{
			"floor 2.9",
			nil,
			block.NewDecimal(decimal.New(2, 0)),
		},
		helpers.Test{
			"floor -2.7",
			nil,
			block.NewDecimal(decimal.New(-3, 0)),
		},
		helpers.Test{
			"floor -2",
			nil,
			block.NewDecimal(decimal.New(-2, 0)),
		},
	)
}

func TestCeil(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"ceil 2",
			nil,
			block.NewDecimal(decimal.New(2, 0)),
		},
		helpers.Test{
			"ceil 2.4",
			nil,
			block.NewDecimal(decimal.New(3, 0)),
		},
		helpers.Test{
			"ceil 2.9",
			nil,
			block.NewDecimal(decimal.New(3, 0)),
		},
		helpers.Test{
			"ceil -2.7",
			nil,
			block.NewDecimal(decimal.New(-2, 0)),
		},
		helpers.Test{
			"ceil -2",
			nil,
			block.NewDecimal(decimal.New(-2, 0)),
		},
	)
}
