package lexer

import (
	"testing"

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
				block.New("1"),
				block.New("2"),
			),
		},
		{
			"+\t1\t2",
			block.New("+",
				block.New("1"),
				block.New("2"),
			),
		},
		{
			"+ (+ 1 2) 3",
			block.New("+",
				block.New("+",
					block.New("1"),
					block.New("2"),
				),
				block.New("3"),
			),
		},
		{
			"+ 1 (- 2 3)",
			block.New("+",
				block.New("1"),
				block.New("-",
					block.New("2"),
					block.New("3"),
				),
			),
		},
		{
			"+ (- 1 2) (* 3 4)",
			block.New("+",
				block.New("-",
					block.New("1"),
					block.New("2"),
				),
				block.New("*",
					block.New("3"),
					block.New("4"),
				),
			),
		},
		{
			"fn A B C",
			block.New("fn",
				block.New("A"),
				block.New("B"),
				block.New("C"),
			),
		},
		{
			`text "Hello World"`,
			block.New("text",
				block.New("Hello World"),
			),
		},
		{
			`text "Hello \"W o r l d\" and Universe"`,
			block.New("text",
				block.New(`Hello "W o r l d" and Universe`),
			),
		},
		{
			"+ 1.2 3.4",
			block.New("+",
				block.New("1.2"),
				block.New("3.4"),
			),
		},
		{
			"+ -1 -2",
			block.New("+",
				block.New("-1"),
				block.New("-2"),
			),
		},
		{
			"+ -1.2 -3.4",
			block.New("+",
				block.New("-1.2"),
				block.New("-3.4"),
			),
		},
		{
			"+ -.1 -.2",
			block.New("+",
				block.New("-.1"),
				block.New("-.2"),
			),
		},
		{
			"(fn A B)",
			block.New("",
				block.New("fn",
					block.New("A"),
					block.New("B"),
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
					block.New("B"),
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
					block.New("1"),
					block.New("1"),
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
					block.New("1"),
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
