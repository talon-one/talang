package block

import (
	"testing"
	"time"

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
	if !ok {
		panic("Not a decimal")
	}
	return d
}
func TestNew(t *testing.T) {
	for _, test := range validDecimalFormats {
		require.EqualValues(t, &Block{Text: test, Kind: DecimalKind, Decimal: mustDecimal(decimal.New(0, 0).SetString(test)), Children: []*Block{}}, New(test), "Failed at %s", test)
	}
	/*for _, test := range invalidDecimalFormats {
		require.EqualValues(t, &Block{Text: test, Kind: StringKind}, New(test), "Failed at %s", test)
	}*/
}

func TestNewTyped(t *testing.T) {
	time, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	if err != nil {
		panic(err)
	}
	tests := []struct {
		expectedKind Kind
		expectedText string
		input        *Block
	}{
		{BoolKind, "false", NewBool(false)},
		{BoolKind, "true", NewBool(true)},

		{StringKind, "Hallo", NewString("Hallo")},

		{DecimalKind, "1", NewDecimal(decimal.New(1, 0))},

		{TimeKind, "2006-01-02T15:04:05Z", NewTime(time)},

		{NullKind, "", NewNull()},
	}
	for _, test := range tests {
		require.Equal(t, test.expectedKind, test.input.Kind)
		require.Equal(t, test.expectedText, test.input.Text)
	}
}

func TestIsDecimal(t *testing.T) {
	block := New(validDecimalFormats[0])
	require.Equal(t, true, block.IsDecimal())
}

func TestIsTime(t *testing.T) {
	block := New("2001-01-01T01:01:01Z")
	require.Equal(t, true, block.IsTime())
}

func TestIsString(t *testing.T) {
	block := New("Hello World")
	require.Equal(t, true, block.IsString())
}

func TestIsBool(t *testing.T) {
	block := New("false")
	require.Equal(t, true, block.IsBool())
	block = New("true")
	require.Equal(t, true, block.IsBool())
}

func TestIsNull(t *testing.T) {
	block := NewNull()
	require.Equal(t, true, block.IsNull())
}

func TestIsEmpty(t *testing.T) {
	var block Block
	require.Equal(t, true, block.IsEmpty())
}

func TestIsBlock(t *testing.T) {
	b := New("Hello", NewString("World"))
	require.Equal(t, true, b.IsBlock())
	b = New("")
	require.Equal(t, true, b.IsBlock())
}

func TestUpdate(t *testing.T) {
	// create a simple block
	var b Block
	b.Update(NewString("Hello"))
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.Text)

	b.Update(NewBool(false))
	require.Equal(t, true, b.IsBool())
	require.Equal(t, "false", b.Text)

	time, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	if err != nil {
		panic(err)
	}

	b.Update(NewTime(time))
	require.Equal(t, true, b.IsTime())
	require.Equal(t, "2006-01-02T15:04:05Z", b.Text)

	b.Update(NewDecimal(decimal.New(1, 0)))
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "1", b.Text)

	b.Update(NewNull())
	require.Equal(t, true, b.IsNull())
	require.Equal(t, "", b.Text)
}

func TestString(t *testing.T) {
	block := New("+", New("1"), New("2"))
	require.Equal(t, "(+ 1 2)", block.String())

	block = New("+", New("-", New("1"), New("2")), New("3"))
	require.Equal(t, "(+ (- 1 2) 3)", block.String())
}

func TestArguments(t *testing.T) {
	block := New("+", New("1"), New("2"))
	require.EqualValues(t, []Kind{DecimalKind, DecimalKind}, Arguments(block.Children))

	block = New("+", New("Hello"), New("1"))
	require.EqualValues(t, []Kind{StringKind, DecimalKind}, Arguments(block.Children))
}

func TestToHumanReadable(t *testing.T) {
	block := New("+", New("1"), New("2"))
	require.Equal(t, "1, 2", BlockArguments(block.Children).ToHumanReadable())
}
