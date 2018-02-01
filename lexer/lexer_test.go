package talang

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/term"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		input    string
		expected term.Term
	}{
		{
			"+ 1 2",
			term.New("+",
				term.New("1"),
				term.New("2"),
			),
		},
		{
			"+ (+ 1 2) 3",
			term.New("+",
				term.New("+",
					term.New("1"),
					term.New("2"),
				),
				term.New("3"),
			),
		},
		{
			"+ 1 (- 2 3)",
			term.New("+",
				term.New("1"),
				term.New("-",
					term.New("2"),
					term.New("3"),
				),
			),
		},
		{
			"+ (- 1 2) (* 3 4)",
			term.New("+",
				term.New("-",
					term.New("1"),
					term.New("2"),
				),
				term.New("*",
					term.New("3"),
					term.New("4"),
				),
			),
		},
		{
			"fn A B C",
			term.New("fn",
				term.New("A"),
				term.New("B"),
				term.New("C"),
			),
		},
		{
			`text "Hello World"`,
			term.New("text",
				term.New("Hello World"),
			),
		},
		{
			`text "Hello \"W o r l d\" and Universe"`,
			term.New("text",
				term.New(`Hello "W o r l d" and Universe`),
			),
		},
		{
			"+ 1.2 3.4",
			term.New("+",
				term.New("1.2"),
				term.New("3.4"),
			),
		},
		{
			"+ -1 -2",
			term.New("+",
				term.New("-1"),
				term.New("-2"),
			),
		},
		{
			"+ -1.2 -3.4",
			term.New("+",
				term.New("-1.2"),
				term.New("-3.4"),
			),
		},
		{
			"+ -.1 -.2",
			term.New("+",
				term.New("-.1"),
				term.New("-.2"),
			),
		},
		{
			"(fn A B)",
			term.New("",
				term.New("fn",
					term.New("A"),
					term.New("B"),
				),
			),
		},
	}

	for i := 0; i < len(tests); i++ {
		s, err := Lex(tests[i].input)
		if err != nil {
			panic(err)
		}
		require.EqualValues(t, tests[i].expected, s, "Input `%s' failed", tests[i].input)
	}

}

func TestUnquote(t *testing.T) {

	doubleQuote := func(extracted string, rest string, str string) {
		a, b := unquote(str, '"', '"', '\\')
		require.Equal(t, extracted, a)
		require.Equal(t, rest, b)
	}

	doubleQuote(`Hello World`, ``, `"Hello World"`)
	doubleQuote(`Hello`, ` World`, `"Hello" World`)
	doubleQuote(`Hello`, ` "World"`, `"Hello" "World"`)
	doubleQuote(`Hello World and Universe`, ``, `"Hello World and Universe"`)
	doubleQuote(`Hello 'World and Universe'`, ``, `"Hello 'World and Universe'"`)

	// nested cases
	doubleQuote(``, `Hello World""`, `""Hello World""`)
	doubleQuote(`Hello "World"`, ``, `"Hello \"World\""`)
	doubleQuote(``, `\"Hello World\"`, `\"Hello World\"`)
	doubleQuote(`Hello "World and \"Universe\""`, ``, `"Hello \"World and \\"Universe\\"\""`)
	doubleQuote(`Hello "World" and Universe`, ` and Dimension`, `"Hello \"World\" and Universe" and Dimension`)

	// invalid cases
	doubleQuote(``, ``, ``)
	doubleQuote(``, `"`, `"`)
	doubleQuote(``, `"Test`, `"Test`)
	doubleQuote(`Hello`, ` "World`, `"Hello" "World`)
	doubleQuote(`Hello "`, ``, `"Hello \""`)
	doubleQuote(``, ``, `""`)

	doubleQuote(``, `Hello World`, `Hello World`)
	doubleQuote(``, `Hello "World"`, `Hello "World"`)

	brakets := func(extracted string, rest string, str string) {
		a, b := unquote(str, '(', ')', utf8.RuneError)
		require.Equal(t, extracted, a)
		require.Equal(t, rest, b)
	}

	brakets(`Hello World`, ` and Universe`, `(Hello World) and Universe`)
	brakets(`echo $(echo Hello World)`, ``, `(echo $(echo Hello World))`)
	brakets(`echo $(echo $(echo Hello World))`, ``, `(echo $(echo $(echo Hello World)))`)
	brakets(`Hello (World (, (Universe,) Dimension)) and Religion`, ``, `(Hello (World (, (Universe,) Dimension)) and Religion)`)

	brakets(``, `(Hello World`, `(Hello World`)
	brakets(`Hello (World`, ``, `(Hello (World)`)
	brakets(`Hello ((World)`, ``, `(Hello ((World))`)

	brakets(`Token2 Token3 (Token4 Token5) Token6`, ` Token7`, `(Token2 Token3 (Token4 Token5) Token6) Token7`)

	brakets(``, `(Token`, `(Token`)

	ticks := func(extracted string, rest string, str string) {
		a, b := unquote(str, '`', '`', utf8.RuneError)
		require.Equal(t, extracted, a)
		require.Equal(t, rest, b)
	}

	ticks("Token2 Token3 Token4 Token5 Token6", " Token7", "`Token2 Token3 Token4 Token5 Token6` Token7")

	curlyBrackets := func(extracted string, rest string, str string) {
		a, b := unquote(str, '{', '}', utf8.RuneError)
		require.Equal(t, extracted, a)
		require.Equal(t, rest, b)
	}

	curlyBrackets(`Hello World`, ` and Universe`, `{Hello World} and Universe`)
	curlyBrackets(`echo ${echo Hello World}`, ``, `{echo ${echo Hello World}}`)
	curlyBrackets(`echo ${echo ${echo Hello World}}`, ``, `{echo ${echo ${echo Hello World}}}`)
	curlyBrackets(`Hello {World {, {Universe,} Dimension}} and Religion`, ``, `{Hello {World {, {Universe,} Dimension}} and Religion}`)

	curlyBrackets(``, `{Hello World`, `{Hello World`)
	curlyBrackets(`Hello {World`, ``, `{Hello {World}`)
	curlyBrackets(`Hello {{World}`, ``, `{Hello {{World}}`)

	curlyBrackets(`Token2 Token3 {Token4 Token5} Token6`, ` Token7`, `{Token2 Token3 {Token4 Token5} Token6} Token7`)

	curlyBrackets(``, `{Token`, `{Token`)

}

func BenchmarkParse(b *testing.B) {
	var tests = []struct {
		input    string
		expected term.Term
	}{
		{
			"(= 1 1)",
			term.New("",
				term.New("=",
					term.New("1"),
					term.New("1"),
				),
			),
		},
		{
			`(= "blah" "bl\"ah")`,
			term.New("",
				term.New("=",
					term.New("blah"),
					term.New("bl\"ah"),
				),
			),
		},
		{
			`(= "what\"" what)`,
			term.New("",
				term.New("=",
					term.New("what\""),
					term.New("what"),
				),
			),
		},
		{
			"(exists (. Session CartItems) (x (= (. x Name) (# 0))))",
			term.New("",
				term.New("exits",
					term.New(".",
						term.New("Session"),
						term.New("CartItems")),
					term.New("x",
						term.New("=",
							term.New(".",
								term.New("x"),
								term.New("Name"),
							),
							term.New("#",
								term.New("0"),
							),
						),
					),
				),
			),
		},
		{
			`(rule (! sessionHasValidCoupon) (! setDiscountAmount "Coupon Discount" 10))`,
			term.New("",
				term.New("rule",
					term.New("!",
						term.New("sessionHasValidCoupon"),
					),
					term.New("!",
						term.New("setDiscountAmount"),
						term.New("Coupon Discount"),
						term.New("10"),
					),
				),
			),
		},
		{
			`(- 1 -1)`,
			term.New("",
				term.New("-",
					term.New("1"),
					term.New("-1"),
				),
			),
		},
		{
			`(-of-course this-is-allowed)`,
			term.New("",
				term.New("-of-course",
					term.New("this-is-alowed"),
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
