package unquote

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnqoute(t *testing.T) {
	tests := []struct {
		Input     string
		Start     string
		End       string
		Extracted string
		Rest      string
	}{
		// Simple Quote
		{
			`(Hello World)`,
			`(`,
			`)`,
			`Hello World`,
			``,
		},
		//Quote in Quote(s)
		{
			`(Hello (World))`,
			`(`,
			`)`,
			`Hello (World)`,
			``,
		},
		{
			`(Hello (World) and Universe)`,
			`(`,
			`)`,
			`Hello (World) and Universe`,
			``,
		},
		{
			`(Hello (World) (Universe))`,
			`(`,
			`)`,
			`Hello (World) (Universe)`,
			``,
		},
		{
			`(Hello (World) and (Universe (and Sun)))`,
			`(`,
			`)`,
			`Hello (World) and (Universe (and Sun))`,
			``,
		},
		{
			`(Hello (World)) (Universe)`,
			`(`,
			`)`,
			`Hello (World)`,
			` (Universe)`,
		},
		{
			`(Hello (World and (Universe))) some Rest`,
			`(`,
			`)`,
			`Hello (World and (Universe))`,
			` some Rest`,
		},
		// Two quotes
		{
			`(Hello) (World)`,
			`(`,
			`)`,
			`Hello`,
			` (World)`,
		},
		// Not starting with a Quote
		{
			`Hello) (World)`,
			`(`,
			`)`,
			``,
			`Hello) (World)`,
		},
		// Not finished quote
		{
			`(Hello`,
			`(`,
			`)`,
			``,
			`(Hello`,
		},
		{
			`(Hello (World)`,
			`(`,
			`)`,
			``,
			`(Hello (World)`,
		},
		// Same Start and End Quote
		{
			`(Hello) (World)`,
			`(`,
			`(`,
			``,
			`(Hello) (World)`,
		},
	}

	for i, test := range tests {
		extracted, rest := Unquote(test.Input, test.Start, test.End)
		require.Equal(t, test.Extracted, extracted, "Test #%d failed (extracted):`%s'", i, test.Input)
		require.Equal(t, test.Rest, rest, "Test #%d failed (rest): `%s'", i, test.Input)
	}
}

func TestEscapeUnquote(t *testing.T) {
	tests := []struct {
		Input     string
		Quote     string
		Escape    string
		Extracted string
		Rest      string
	}{
		// Simple Quote
		{
			`"Hello World"`,
			`"`,
			`\`,
			`Hello World`,
			``,
		},
		//Quote in Quote(s)
		{
			`"Hello \"World\""`,
			`"`,
			`\`,
			`Hello "World"`,
			``,
		},
		{
			`"Hello \"World\" and Universe"`,
			`"`,
			`\`,
			`Hello "World" and Universe`,
			``,
		},
		{
			`"Hello \"World\" \"Universe\""`,
			`"`,
			`\`,
			`Hello "World" "Universe"`,
			``,
		},
		{
			`"Hello \"World \\"Foo\\"\""`,
			`"`,
			`\`,
			`Hello "World \"Foo\""`,
			``,
		},
		{
			`"Hello \"World \\"Foo \\\"Bar\\\"\\"\""`,
			`"`,
			`\`,
			`Hello "World \"Foo \\"Bar\\"\""`,
			``,
		},
		// Two quotes
		{
			`"Hello" "World"`,
			`"`,
			`\`,
			`Hello`,
			` "World"`,
		},
		{
			`"Hello" \"World\"`,
			`"`,
			`\`,
			`Hello`,
			` \"World\"`,
		},
		// Not starting with a Quote
		{
			`Hello" "World"`,
			`"`,
			`\`,
			``,
			`Hello" "World"`,
		},
		// Not finished quote
		{
			`"Hello`,
			`"`,
			`\`,
			``,
			`"Hello`,
		},
		{
			`"Hello \"World\"`,
			`"`,
			`\`,
			``,
			`"Hello \"World\"`,
		},
		{
			`"Hello \"World"`,
			`"`,
			`\`,
			`Hello "World`,
			``,
		},
	}

	for i, test := range tests {
		extracted, rest := EscapeUnquote(test.Input, test.Quote, test.Escape)
		require.Equal(t, test.Extracted, extracted, "Test #%d failed (extracted):`%s'", i, test.Input)
		require.Equal(t, test.Rest, rest, "Test #%d failed (rest): `%s'", i, test.Input)
	}
}

// Some legacy tests, keep them in, just to be sure
func TestUnquoteLegacy(t *testing.T) {
	doubleQuote := func(extracted string, rest string, str string) {
		a, b := EscapeUnquote(str, `"`, `\`)
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
		a, b := Unquote(str, `(`, `)`)
		require.Equal(t, extracted, a)
		require.Equal(t, rest, b)
	}

	brakets(`Hello World`, ` and Universe`, `(Hello World) and Universe`)
	brakets(`echo $(echo Hello World)`, ``, `(echo $(echo Hello World))`)
	brakets(`echo $(echo $(echo Hello World))`, ``, `(echo $(echo $(echo Hello World)))`)
	brakets(`Hello (World (, (Universe,) Dimension)) and Religion`, ``, `(Hello (World (, (Universe,) Dimension)) and Religion)`)
	brakets(`A (B ((C) D)) E F`, ``, `(A (B ((C) D)) E F)`)
	brakets(`(C) (D)`, ``, `((C) (D))`)

	brakets(``, `(Hello World`, `(Hello World`)
	brakets(``, `(Hello (World)`, `(Hello (World)`)
	brakets(``, `(Hello ((World))`, `(Hello ((World))`)

	brakets(`Token2 Token3 (Token4 Token5) Token6`, ` Token7`, `(Token2 Token3 (Token4 Token5) Token6) Token7`)
	brakets(`Token1 (Token2 Token3) (Token4 (Token5 Token6))`, ``, `(Token1 (Token2 Token3) (Token4 (Token5 Token6)))`)

	brakets(``, `(Token`, `(Token`)
	brakets(`Token`, ` `, `(Token) `)
	brakets(``, ` (Token)`, ` (Token)`)

	ticks := func(extracted string, rest string, str string) {
		a, b := EscapeUnquote(str, "`", `\`)
		require.Equal(t, extracted, a)
		require.Equal(t, rest, b)
	}

	ticks("Token2 Token3 Token4 Token5 Token6", " Token7", "`Token2 Token3 Token4 Token5 Token6` Token7")

	curlyBrackets := func(extracted string, rest string, str string) {
		a, b := Unquote(str, `{`, `}`)
		require.Equal(t, extracted, a)
		require.Equal(t, rest, b)
	}

	curlyBrackets(`Hello World`, ` and Universe`, `{Hello World} and Universe`)
	curlyBrackets(`echo ${echo Hello World}`, ``, `{echo ${echo Hello World}}`)
	curlyBrackets(`echo ${echo ${echo Hello World}}`, ``, `{echo ${echo ${echo Hello World}}}`)
	curlyBrackets(`Hello {World {, {Universe,} Dimension}} and Religion`, ``, `{Hello {World {, {Universe,} Dimension}} and Religion}`)

	curlyBrackets(``, `{Hello World`, `{Hello World`)
	curlyBrackets(``, `{Hello {World}`, `{Hello {World}`)
	curlyBrackets(``, `{Hello {{World}}`, `{Hello {{World}}`)

	curlyBrackets(`Token2 Token3 {Token4 Token5} Token6`, ` Token7`, `{Token2 Token3 {Token4 Token5} Token6} Token7`)

	curlyBrackets(``, `{Token`, `{Token`)
}
