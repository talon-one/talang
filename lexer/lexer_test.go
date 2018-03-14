package lexer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected *block.Block
	}{
		{
			"+ 1 2",
			block.New("+",
				block.NewDecimalFromInt(1),
				block.NewDecimalFromInt(2),
			),
		},
		{
			"+\t1\t2",
			block.New("+",
				block.NewDecimalFromInt(1),
				block.NewDecimalFromInt(2),
			),
		},
		{
			"+ (+ 1 2) 3",
			block.New("+",
				block.New("+",
					block.NewDecimalFromInt(1),
					block.NewDecimalFromInt(2),
				),
				block.NewDecimalFromInt(3),
			),
		},
		{
			"+ 1 (- 2 3)",
			block.New("+",
				block.NewDecimalFromInt(1),
				block.New("-",
					block.NewDecimalFromInt(2),
					block.NewDecimalFromInt(3),
				),
			),
		},
		{
			"+ (- 1 2) (* 3 4)",
			block.New("+",
				block.New("-",
					block.NewDecimalFromInt(1),
					block.NewDecimalFromInt(2),
				),
				block.New("*",
					block.NewDecimalFromInt(3),
					block.NewDecimalFromInt(4),
				),
			),
		},
		{
			"fn A B C",
			block.New("fn",
				block.NewString("A"),
				block.NewString("B"),
				block.NewString("C"),
			),
		},
		{
			`text "Hello World"`,
			block.New("text",
				block.NewString("Hello World"),
			),
		},
		{
			`text "Hello \"W o r l d\" and Universe"`,
			block.New("text",
				block.NewString(`Hello "W o r l d" and Universe`),
			),
		},
		{
			"+ 1.2 3.4",
			block.New("+",
				block.NewDecimalFromString("1.2"),
				block.NewDecimalFromString("3.4"),
			),
		},
		{
			"+ -1 -2",
			block.New("+",
				block.NewDecimalFromInt(-1),
				block.NewDecimalFromInt(-2),
			),
		},
		{
			"+ -1.2 -3.4",
			block.New("+",
				block.NewDecimalFromString("-1.2"),
				block.NewDecimalFromString("-3.4"),
			),
		},
		{
			"+ -.1 -.2",
			block.New("+",
				block.NewDecimalFromString("-.1"),
				block.NewDecimalFromString("-.2"),
			),
		},
		{
			"(fn A B)",
			block.New("",
				block.New("fn",
					block.NewString("A"),
					block.NewString("B"),
				),
			),
		},
		{
			"(fn (A) B)",
			block.New("",
				block.New("fn",
					&block.Block{
						String:   "A",
						Kind:     block.BlockKind,
						Children: []*block.Block{},
					},
					block.NewString("B"),
				),
			),
		},
		{
			`set CartItems (push (. CartItems))`,
			block.New("set",
				block.NewString("CartItems"),
				block.New("push",
					block.New(".",
						block.NewString("CartItems"),
					),
				),
			),
		},
		{
			`set CartItems (push (. CartItems) (kv (Position 99)))`,
			block.New("set",
				block.NewString("CartItems"),
				block.New("push",
					block.New(".", block.NewString("CartItems")),
					block.New("kv",
						block.New("Position", block.NewDecimalFromInt(99)),
					),
				),
			),
		},
	}

	for i := 0; i < len(tests); i++ {
		s, err := Lex(tests[i].input)
		if err != nil {
			panic(err)
		}
		require.Equal(t, true, tests[i].expected.Equal(s), "Input `%s' failed, was `%s'", tests[i].input, s.Stringify())
	}
}

func TestForcedString(t *testing.T) {
	timeMustParse := func(s string) time.Time {
		tm, err := time.Parse(time.RFC3339, s)
		if err != nil {
			panic(err)
		}
		return tm
	}

	tests := []struct {
		input    string
		expected *block.Block
	}{
		{
			`+ "1" "2"`,
			block.New("+",
				block.NewString("1"),
				block.NewString("2"),
			),
		},
		{
			`+ 1 2`,
			block.New("+",
				block.NewDecimalFromInt(1),
				block.NewDecimalFromInt(2),
			),
		},
		{
			`+ "true" "false"`,
			block.New("+",
				block.NewString("true"),
				block.NewString("false"),
			),
		},
		{
			`+ true false`,
			block.New("+",
				block.NewBool(true),
				block.NewBool(false),
			),
		},
		{
			`+ "2007-01-02T00:00:00Z" "2007-01-02T00:00:00Z"`,
			block.New("+",
				block.NewString("2007-01-02T00:00:00Z"),
				block.NewString("2007-01-02T00:00:00Z"),
			),
		},
		{
			`+ 2007-01-02T00:00:00Z 2007-01-02T00:00:00Z`,
			block.New("+",
				block.NewTime(timeMustParse("2007-01-02T00:00:00Z")),
				block.NewTime(timeMustParse("2007-01-02T00:00:00Z")),
			),
		},
	}
	for i := 0; i < len(tests); i++ {
		s, err := Lex(tests[i].input)
		if err != nil {
			panic(err)
		}
		require.Equal(t, true, tests[i].expected.Equal(s), "Input `%s' failed, was `%s'", tests[i].input, s.Stringify())
	}
}

func TestMustLex(t *testing.T) {
	require.NotPanics(t, func() {
		MustLex("(Hello)").Equal(block.NewString("Hello"))
	})
}

func BenchmarkParse(b *testing.B) {
	var tests = []struct {
		input    string
		expected *block.Block
	}{
		{
			"(= 1 1)",
			block.New("",
				block.New("=",
					block.NewDecimalFromInt(1),
					block.NewDecimalFromInt(1),
				),
			),
		},
		{
			`(= "blah" "bl\"ah")`,
			block.New("",
				block.New("=",
					block.New("blah"),
					block.New("bl\"ah"),
				),
			),
		},
		{
			`(= "what\"" what)`,
			block.New("",
				block.New("=",
					block.New("what\""),
					block.New("what"),
				),
			),
		},
		{
			"(exists (. Session CartItems) (x (= (. x Name) (# 0))))",
			block.New("",
				block.New("exits",
					block.New(".",
						block.New("Session"),
						block.New("CartItems")),
					block.New("x",
						block.New("=",
							block.New(".",
								block.New("x"),
								block.New("Name"),
							),
							block.New("#",
								block.New("0"),
							),
						),
					),
				),
			),
		},
		{
			`(rule (! sessionHasValidCoupon) (! setDiscountAmount "Coupon Discount" 10))`,
			block.New("",
				block.New("rule",
					block.New("!",
						block.New("sessionHasValidCoupon"),
					),
					block.New("!",
						block.New("setDiscountAmount"),
						block.New("Coupon Discount"),
						block.New("10"),
					),
				),
			),
		},
		{
			`(- 1 -1)`,
			block.New("",
				block.New("-",
					block.NewDecimalFromInt(1),
					block.New("-1"),
				),
			),
		},
		{
			`(-of-course this-is-allowed)`,
			block.New("",
				block.New("-of-course",
					block.New("this-is-allowed"),
				),
			),
		},
	}

	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			_, err := Lex(test.input)
			if err != nil {
				panic(err)
			}
		}
	}
}
