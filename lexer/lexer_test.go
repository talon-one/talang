package lexer

import (
	"testing"
	"unicode/utf8"

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
						Text:     "A",
						Kind:     block.BlockKind,
						Children: []*block.Block{},
					},
					block.New("B"),
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
