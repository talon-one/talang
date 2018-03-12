package block

import (
	"encoding/json"
	"sort"
	"testing"
	"time"

	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/require"
)

func mustDecimal(d *decimal.Big, ok bool) *decimal.Big {
	if !ok {
		panic("Not a decimal")
	}
	return d
}
func TestNew(t *testing.T) {
	tests := []struct {
		input    string
		expected *Block
	}{
		// Valid Decimal Formats
		{
			"1",
			NewDecimalFromInt(1),
		},
		{
			"+1",
			NewDecimalFromInt(1),
		},
		{
			"-1",
			NewDecimalFromInt(-1),
		},
		{
			"1.2",
			NewDecimalFromString("1.2"),
		},
		{
			"+1.2",
			NewDecimalFromString("1.2"),
		},
		{
			"-1.2",
			NewDecimalFromString("-1.2"),
		},
		{
			".1",
			NewDecimalFromString("0.1"),
		},
		{
			"+.1",
			NewDecimalFromString("0.1"),
		},
		{
			"-.1",
			NewDecimalFromString("-0.1"),
		},
		{
			"1.",
			NewDecimalFromInt(1),
		},
		{
			"+1.",
			NewDecimalFromInt(1),
		},
		{
			"-1.",
			NewDecimalFromInt(-1),
		},
		// {
		// 	"1e2",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"+1e2",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"-1e2",
		// 	NewDecimalFromInt(-1),
		// },
		// {
		// 	"1e+2",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"1e-2",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"+1e+2",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"-1e-2",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"+1e-2",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"-1e+2",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"1.2e3",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"+1.2e3",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"-1.2e3",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"1.2e+3",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"1.2e-3",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"+1.2e-3",
		// 	NewDecimalFromInt(1),
		// },
		// {
		// 	"-1.2e+3",
		// 	NewDecimalFromInt(1),
		// },
		{
			"Hello",
			NewString("Hello"),
		},
		{
			"1+2",
			NewString("1+2"),
		},
		{
			"1-2",
			NewString("1-2"),
		},
		{
			"1.+2",
			NewString("1.+2"),
		},
		{
			"1.-2",
			NewString("1.-2"),
		},
		{
			"++1",
			NewString("++1"),
		},
		{
			"--1",
			NewString("--1"),
		},
		{
			"1+",
			NewString("1+"),
		},
		{
			"1-",
			NewString("1-"),
		},
		{
			"1e",
			NewString("1e"),
		},
		{
			"+1e",
			NewString("+1e"),
		},
		{
			"-1e",
			NewString("-1e"),
		},
		{
			"e1",
			NewString("e1"),
		},
		{
			"e+1",
			NewString("e+1"),
		},
		{
			"e-1",
			NewString("e-1"),
		},
		{
			"1.2.3",
			NewString("1.2.3"),
		},
	}
	for i, test := range tests {
		require.Equal(t, true, test.expected.Equal(New(test.input)), "Test %d failed: `%s'", i, test.input)
	}
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
		{DecimalKind, "2", NewDecimalFromInt(2)},
		{DecimalKind, "3", NewDecimalFromString("3")},
		{NullKind, "", NewDecimalFromString("HELLO3HELLO")},

		{TimeKind, "2006-01-02T15:04:05Z", NewTime(time)},

		{ListKind, "", NewList()},

		{NullKind, "", NewNull()},
	}
	for _, test := range tests {
		require.Equal(t, test.expectedKind, test.input.Kind)
		require.Equal(t, test.expectedText, test.input.String)
	}
}

func TestIsDecimal(t *testing.T) {
	block := New("1")
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
	require.Equal(t, "Hello", b.String)

	b.Update(NewBool(false))
	require.Equal(t, true, b.IsBool())
	require.Equal(t, "false", b.String)

	time, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	if err != nil {
		panic(err)
	}

	b.Update(NewTime(time))
	require.Equal(t, true, b.IsTime())
	require.Equal(t, "2006-01-02T15:04:05Z", b.String)

	b.Update(NewDecimal(decimal.New(1, 0)))
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "1", b.String)

	b.Update(NewNull())
	require.Equal(t, true, b.IsNull())
	require.Equal(t, "", b.String)

	b.Update(NewList(NewString("Hello"), NewString("World")))
	require.Equal(t, true, b.IsList())
	require.EqualValues(t, NewList(NewString("Hello"), NewString("World")).Children, b.Children)

	b.Update(NewMap(map[string]*Block{
		"Key1": NewString("Value1"),
		"Key2": NewList(NewString("Value2"), NewString("Value3")),
	}))
	require.Equal(t, true, b.IsMap())
	require.EqualValues(t, map[string]*Block{
		"Key1": NewString("Value1"),
		"Key2": NewList(NewString("Value2"), NewString("Value3")),
	}, b.Map())
}

func TestMapItem(t *testing.T) {
	block := NewMap(map[string]*Block{
		"Key1": NewBool(true),
		"Key2": NewDecimalFromInt(1),
		"Key3": NewString("Hello"),
	})
	require.Equal(t, NewBool(true), block.MapItem("Key1"))

	require.Equal(t, NewNull(), block.MapItem("Key4"))
}

func TestSetMapItem(t *testing.T) {
	block := NewMap(map[string]*Block{
		"Key1": NewBool(true),
		"Key2": NewDecimalFromInt(1),
		"Key3": NewString("Hello"),
	})
	block.SetMapItem("Key1", NewString("World"))
	require.Equal(t, NewString("World"), block.MapItem("Key1"))

	block.SetMapItem("Key4", NewString("Foo"))
	require.Equal(t, NewString("Foo"), block.MapItem("Key4"))
}

func TestStringify(t *testing.T) {
	block := New("+", New("1"), New("2"))
	require.Equal(t, "(+ 1 2)", block.Stringify())

	block = New("+", New("-", New("1"), New("2")), New("3"))
	require.Equal(t, "(+ (- 1 2) 3)", block.Stringify())

	block = New("", New("-", New("1"), New("2")), New("3"))
	require.Equal(t, "((- 1 2) 3)", block.Stringify())
}

func TestEqual(t *testing.T) {
	now := time.Now()
	tests := []struct {
		a        *Block
		b        *Block
		expected bool
	}{
		{
			NewDecimalFromInt(1),
			NewDecimalFromInt(1),
			true,
		},
		{
			NewString("Hello"),
			NewString("Hello"),
			true,
		},
		{
			NewBool(true),
			NewBool(true),
			true,
		},
		{
			NewTime(now),
			NewTime(now),
			true,
		},
		{
			NewNull(),
			NewNull(),
			true,
		},
		{
			NewList(NewBool(true), NewString("Hello")),
			NewList(NewBool(true), NewString("Hello")),
			true,
		},
		{
			NewMap(map[string]*Block{
				"Key1": NewBool(true),
				"Key2": NewDecimalFromInt(1),
				"Key3": NewString("Hello"),
			}),
			NewMap(map[string]*Block{
				"Key2": NewDecimalFromInt(1),
				"Key1": NewBool(true),
				"Key3": NewString("Hello"),
			}),
			true,
		},
		{
			New("", NewBool(true), NewString("Hello")),
			New("", NewBool(true), NewString("Hello")),
			true,
		},
		{
			NewDecimalFromInt(1),
			nil,
			false,
		},
		{
			&Block{},
			&Block{},
			false,
		},
		{
			NewMap(map[string]*Block{
				"Key1": NewBool(true),
				"Key3": NewString("Hello"),
			}),
			NewMap(map[string]*Block{
				"Key2": NewDecimalFromInt(1),
				"Key1": NewBool(true),
				"Key3": NewString("Hello"),
			}),
			false,
		},
		{
			NewMap(map[string]*Block{
				"Key1": NewBool(true),
				"Key2": NewDecimalFromInt(2),
				"Key3": NewString("Hello"),
			}),
			NewMap(map[string]*Block{
				"Key2": NewDecimalFromInt(1),
				"Key1": NewBool(true),
				"Key3": NewString("Hello"),
			}),
			false,
		},
		{
			New("", NewBool(true)),
			New("", NewBool(true), NewString("Hello")),
			false,
		},
		{
			New("", NewBool(true), NewString("Hello")),
			New("", NewBool(false), NewString("Hello")),
			false,
		},
	}
	for i, test := range tests {
		require.Equal(t, test.expected, test.a.Equal(test.b), "Test %d failed", i)
	}
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

func TestMarshaling(t *testing.T) {
	block1 := NewMap(map[string]*Block{
		"Key2": NewDecimalFromInt(1),
		"Key1": NewBool(true),
		"Key3": NewString("Hello"),
		"Key4": NewList(NewBool(false), NewMap(map[string]*Block{
			"SubKey1": NewDecimalFromInt(3),
			"SubKey2": NewTime(time.Now()),
		})),
	})
	b, err := json.Marshal(block1)
	require.NoError(t, err)

	var block2 Block

	require.NoError(t, json.Unmarshal(b, &block2))

	require.Equal(t, true, block1.Equal(&block2))
}

func TestSort(t *testing.T) {
	list := NewList(NewDecimalFromInt(4), NewDecimalFromInt(2), NewDecimalFromInt(3))
	sort.Sort(BlockArguments(list.Children))
	require.Equal(t, true, list.Equal(NewList(NewDecimalFromInt(2), NewDecimalFromInt(3), NewDecimalFromInt(4))))
}
