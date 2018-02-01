package term

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

// Valid
var validFormats = []string{
	"1",
	"+1",
	"-1",
	"1.2",
	"+1.2",
	"-1.2",
	".1",
	"+.1",
	"-.1",
	"1.",
	"+1.",
	"-1.",
	"1e2",
	"+1e2",
	"-1e2",
	"1e+2",
	"1e-2",
	"+1e+2",
	"-1e-2",
	"+1e-2",
	"-1e+2",
	"1.2e3",
	"+1.2e3",
	"-1.2e3",
	"1.2e+3",
	"1.2e-3",
	"+1.2e-3",
	"-1.2e+3",
}

// Invalid
var invalidFormats = []string{
	"Hello",
	"1+2",
	"1-2",
	"1.+2",
	"1.-2",
	"++1",
	"--1",
	"1+",
	"1-",
	"1e",
	"+1e",
	"-1e",
	"e1",
	"e+1",
	"e-1",
	"1.2.3",
}

func mustDecimal(d decimal.Decimal, err error) decimal.Decimal {
	if err != nil {
		panic(err)
	}
	return d
}
func TestNew(t *testing.T) {
	for _, test := range validFormats {
		require.EqualValues(t, Term{Text: test, isDecimal: true, Decimal: mustDecimal(decimal.NewFromString(test))}, New(test))
	}
	for _, test := range invalidFormats {
		require.EqualValues(t, Term{Text: test, isDecimal: false}, New(test))
	}
}

func TestIsDecimal(t *testing.T) {
	term := New(validFormats[0])
	require.Equal(t, term.isDecimal, term.IsDecimal())
}

func TestIsEmpty(t *testing.T) {
	var term Term
	require.Equal(t, true, term.IsEmpty())
}

func TestUpdate(t *testing.T) {
	// create a simple term
	term := New(validFormats[0])
	require.Equal(t, validFormats[0], term.Text)
	require.Equal(t, true, term.isDecimal)
	require.Equal(t, 0, len(term.Children))

	// update the text
	term.Update(validFormats[1])
	require.Equal(t, validFormats[1], term.Text)
	require.Equal(t, true, term.isDecimal)
	require.Equal(t, 0, len(term.Children))

	// update the text and their children
	term.Update(validFormats[2], New(validFormats[1]))
	require.Equal(t, validFormats[2], term.Text)
	require.Equal(t, true, term.isDecimal)
	require.Equal(t, 1, len(term.Children))

	// update to an invalid format
	term.Update(invalidFormats[0])
	require.Equal(t, invalidFormats[0], term.Text)
	require.Equal(t, false, term.isDecimal)
	require.Equal(t, 0, len(term.Children))
}

func TestString(t *testing.T) {
	term := New("+", New("1"), New("2"))
	require.Equal(t, "(+ 1 2)", term.String())

	term = New("+", New("-", New("1"), New("2")), New("3"))
	require.Equal(t, "(+ (- 1 2) 3)", term.String())
}
