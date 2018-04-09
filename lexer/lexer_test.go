package lexer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/token"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected *token.TaToken
	}{
		{
			"1",
			token.NewDecimalFromInt(1),
		},
		{
			"+ 1 2",
			token.New("+",
				token.NewDecimalFromInt(1),
				token.NewDecimalFromInt(2),
			),
		},
		{
			"+\t1\t2",
			token.New("+",
				token.NewDecimalFromInt(1),
				token.NewDecimalFromInt(2),
			),
		},
		{
			"+ (+ 1 2) 3",
			token.New("+",
				token.New("+",
					token.NewDecimalFromInt(1),
					token.NewDecimalFromInt(2),
				),
				token.NewDecimalFromInt(3),
			),
		},
		{
			"+ 1 (- 2 3)",
			token.New("+",
				token.NewDecimalFromInt(1),
				token.New("-",
					token.NewDecimalFromInt(2),
					token.NewDecimalFromInt(3),
				),
			),
		},
		{
			"+ (- 1 2) (* 3 4)",
			token.New("+",
				token.New("-",
					token.NewDecimalFromInt(1),
					token.NewDecimalFromInt(2),
				),
				token.New("*",
					token.NewDecimalFromInt(3),
					token.NewDecimalFromInt(4),
				),
			),
		},
		{
			"fn A B C",
			token.New("fn",
				token.NewString("A"),
				token.NewString("B"),
				token.NewString("C"),
			),
		},
		{
			`text "Hello World"`,
			token.New("text",
				token.NewString("Hello World"),
			),
		},
		{
			`text "Hello \"W o r l d\" and Universe"`,
			token.New("text",
				token.NewString(`Hello "W o r l d" and Universe`),
			),
		},
		{
			"+ 1.2 3.4",
			token.New("+",
				token.NewDecimalFromString("1.2"),
				token.NewDecimalFromString("3.4"),
			),
		},
		{
			"+ -1 -2",
			token.New("+",
				token.NewDecimalFromInt(-1),
				token.NewDecimalFromInt(-2),
			),
		},
		{
			"+ -1.2 -3.4",
			token.New("+",
				token.NewDecimalFromString("-1.2"),
				token.NewDecimalFromString("-3.4"),
			),
		},
		{
			"+ -.1 -.2",
			token.New("+",
				token.NewDecimalFromString("-.1"),
				token.NewDecimalFromString("-.2"),
			),
		},
		{
			"(fn A B)",
			token.New("fn",
				token.NewString("A"),
				token.NewString("B"),
			),
		},
		{
			"(fn (A) B)",
			token.New("fn",
				&token.TaToken{
					String:   "A",
					Kind:     token.Token,
					Children: []*token.TaToken{},
				},
				token.NewString("B"),
			),
		},
		{
			`set CartItems (push (. CartItems))`,
			token.New("set",
				token.NewString("CartItems"),
				token.New("push",
					token.New(".",
						token.NewString("CartItems"),
					),
				),
			),
		},
		{
			`set CartItems (push (. CartItems) (kv (Position 99)))`,
			token.New("set",
				token.NewString("CartItems"),
				token.New("push",
					token.New(".", token.NewString("CartItems")),
					token.New("kv",
						token.New("Position", token.NewDecimalFromInt(99)),
					),
				),
			),
		},
		{
			`((+ A B) (+ C D))`,
			token.New("",
				token.New("+",
					token.NewString("A"),
					token.NewString("B"),
				),
				token.New("+",
					token.NewString("C"),
					token.NewString("D"),
				),
			),
		},
	}

	for i, test := range tests {
		s, err := Lex(test.input)
		if err != nil {
			panic(err)
		}
		require.Equal(t, true, test.expected.Equal(s), "Test %d (`%s') failed, was `%s'", i, test.input, s.Stringify())
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
		expected *token.TaToken
	}{
		{
			`+ "1" "2"`,
			token.New("+",
				token.NewString("1"),
				token.NewString("2"),
			),
		},
		{
			`+ 1 2`,
			token.New("+",
				token.NewDecimalFromInt(1),
				token.NewDecimalFromInt(2),
			),
		},
		{
			`+ "true" "false"`,
			token.New("+",
				token.NewString("true"),
				token.NewString("false"),
			),
		},
		{
			`+ true false`,
			token.New("+",
				token.NewBool(true),
				token.NewBool(false),
			),
		},
		{
			`+ "2007-01-02T00:00:00Z" "2007-01-02T00:00:00Z"`,
			token.New("+",
				token.NewString("2007-01-02T00:00:00Z"),
				token.NewString("2007-01-02T00:00:00Z"),
			),
		},
		{
			`+ 2007-01-02T00:00:00Z 2007-01-02T00:00:00Z`,
			token.New("+",
				token.NewTime(timeMustParse("2007-01-02T00:00:00Z")),
				token.NewTime(timeMustParse("2007-01-02T00:00:00Z")),
			),
		},
	}
	for i, test := range tests {
		s, err := Lex(test.input)
		if err != nil {
			panic(err)
		}
		require.Equal(t, true, test.expected.Equal(s), "Test %d (`%s') failed, was `%s'", i, test.input, s.Stringify())
	}
}

func TestMustLex(t *testing.T) {
	require.NotPanics(t, func() {
		MustLex("(Hello)").Equal(token.NewString("Hello"))
	})
}

func BenchmarkParse(b *testing.B) {
	var tests = []struct {
		input    string
		expected *token.TaToken
	}{
		{
			"(= 1 1)",
			token.New("",
				token.New("=",
					token.NewDecimalFromInt(1),
					token.NewDecimalFromInt(1),
				),
			),
		},
		{
			`(= "blah" "bl\"ah")`,
			token.New("",
				token.New("=",
					token.New("blah"),
					token.New("bl\"ah"),
				),
			),
		},
		{
			`(= "what\"" what)`,
			token.New("",
				token.New("=",
					token.New("what\""),
					token.New("what"),
				),
			),
		},
		{
			"(exists (. Session CartItems) (x (= (. x Name) (# 0))))",
			token.New("",
				token.New("exits",
					token.New(".",
						token.New("Session"),
						token.New("CartItems")),
					token.New("x",
						token.New("=",
							token.New(".",
								token.New("x"),
								token.New("Name"),
							),
							token.New("#",
								token.New("0"),
							),
						),
					),
				),
			),
		},
		{
			`(rule (! sessionHasValidCoupon) (! setDiscountAmount "Coupon Discount" 10))`,
			token.New("",
				token.New("rule",
					token.New("!",
						token.New("sessionHasValidCoupon"),
					),
					token.New("!",
						token.New("setDiscountAmount"),
						token.New("Coupon Discount"),
						token.New("10"),
					),
				),
			),
		},
		{
			`(- 1 -1)`,
			token.New("",
				token.New("-",
					token.NewDecimalFromInt(1),
					token.New("-1"),
				),
			),
		},
		{
			`(-of-course this-is-allowed)`,
			token.New("",
				token.New("-of-course",
					token.New("this-is-allowed"),
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
