package block

import (
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/require"
)

// Valid
var validDecimalFormats = []string{
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
var invalidDecimalFormats = []string{
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

func mustDecimal(d *decimal.Big, ok bool) *decimal.Big {
	if ok != true {
		panic("Not a decimal")
	}
	return d
}
func TestNew(t *testing.T) {
	for _, test := range validDecimalFormats {
		require.EqualValues(t, &Block{Text: test, Kind: DecimalKind, Decimal: mustDecimal(decimal.New(0, 0).SetString(test))}, New(test), "Failed at %s", test)
	}
	/*for _, test := range invalidDecimalFormats {
		require.EqualValues(t, &Block{Text: test, Kind: StringKind}, New(test), "Failed at %s", test)
	}*/
}

func TestIsDecimal(t *testing.T) {
	block := New(validDecimalFormats[0])
	require.Equal(t, true, block.IsDecimal())
}

func TestIsEmpty(t *testing.T) {
	var block Block
	require.Equal(t, true, block.IsEmpty())
}

// func TestUpdate(t *testing.T) {
// 	// create a simple block
// 	block := New(validDecimalFormats[0])
// 	require.Equal(t, validDecimalFormats[0], block.Text)
// 	require.Equal(t, true, block.isDecimal)
// 	require.Equal(t, 0, len(block.Children))

// 	// update the text
// 	block.Update(validDecimalFormats[1])
// 	require.Equal(t, validDecimalFormats[1], block.Text)
// 	require.Equal(t, true, block.isDecimal)
// 	require.Equal(t, 0, len(block.Children))

// 	// update the text and their children
// 	block.Update(validDecimalFormats[2], New(validDecimalFormats[1]))
// 	require.Equal(t, validDecimalFormats[2], block.Text)
// 	require.Equal(t, true, block.isDecimal)
// 	require.Equal(t, 1, len(block.Children))

// 	// update to an invalid format
// 	block.Update(invalidDecimalFormats[0])
// 	require.Equal(t, invalidDecimalFormats[0], block.Text)
// 	require.Equal(t, false, block.isDecimal)
// 	require.Equal(t, 0, len(block.Children))
// }

func TestString(t *testing.T) {
	block := New("+", New("1"), New("2"))
	require.Equal(t, "(+ 1 2)", block.String())

	block = New("+", New("-", New("1"), New("2")), New("3"))
	require.Equal(t, "(+ (- 1 2) 3)", block.String())
}

func TestArguments(t *testing.T) {
	block := New("+", New("1"), New("2"))
	require.EqualValues(t, []Kind{DecimalKind, DecimalKind}, block.Arguments())

	block = New("+", New("Hello"), New("1"))
	require.EqualValues(t, []Kind{StringKind, DecimalKind}, block.Arguments())
}
