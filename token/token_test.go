package token

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
		expected *TaToken
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
		input        *TaToken
	}{
		{Boolean, "false", NewBool(false)},
		{Boolean, "true", NewBool(true)},

		{String, "Hallo", NewString("Hallo")},

		{Decimal, "1", NewDecimal(decimal.New(1, 0))},
		{Decimal, "2", NewDecimalFromInt(2)},
		{Decimal, "3", NewDecimalFromString("3")},
		{Null, "", NewDecimalFromString("HELLO3HELLO")},

		{Time, "2006-01-02T15:04:05Z", NewTime(time)},

		{List, "", NewList()},

		{Null, "", NewNull()},
	}
	for _, test := range tests {
		require.Equal(t, test.expectedKind, test.input.Kind)
		require.Equal(t, test.expectedText, test.input.String)
	}
}

func TestNewTime(t *testing.T) {
	parsedTime, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	require.NoError(t, err)
	require.Equal(t, NewTime(parsedTime), New("2006-01-02T15:04:05Z"))
}
func TestIsDecimal(t *testing.T) {
	tkn := New("1")
	require.Equal(t, true, tkn.IsDecimal())
}

func TestIsTime(t *testing.T) {
	tkn := New("2001-01-01T01:01:01Z")
	require.Equal(t, true, tkn.IsTime())
}

func TestIsString(t *testing.T) {
	tkn := New("Hello World")
	require.Equal(t, true, tkn.IsString())
}

func TestIsBool(t *testing.T) {
	tkn := New("false")
	require.Equal(t, true, tkn.IsBool())
	tkn = New("true")
	require.Equal(t, true, tkn.IsBool())
}

func TestIsNull(t *testing.T) {
	tkn := NewNull()
	require.Equal(t, true, tkn.IsNull())
}

func TestIsEmpty(t *testing.T) {
	var tkn TaToken
	require.Equal(t, true, tkn.IsEmpty())
}

func TestIsBlock(t *testing.T) {
	b := New("Hello", NewString("World"))
	require.Equal(t, true, b.IsBlock())
	b = New("")
	require.Equal(t, true, b.IsBlock())
}

func TestCopy(t *testing.T) {
	// create a simple block
	var b TaToken
	Copy(&b, NewString("Hello"))
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.String)

	Copy(&b, NewBool(false))
	require.Equal(t, true, b.IsBool())
	require.Equal(t, "false", b.String)

	time, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	if err != nil {
		panic(err)
	}

	Copy(&b, NewTime(time))
	require.Equal(t, true, b.IsTime())
	require.Equal(t, "2006-01-02T15:04:05Z", b.String)

	Copy(&b, NewDecimal(decimal.New(1, 0)))
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "1", b.String)

	Copy(&b, NewNull())
	require.Equal(t, true, b.IsNull())
	require.Equal(t, "", b.String)

	Copy(&b, NewList(NewString("Hello"), NewString("World")))
	require.Equal(t, true, b.IsList())
	require.EqualValues(t, NewList(NewString("Hello"), NewString("World")).Children, b.Children)

	Copy(&b, NewMap(map[string]*TaToken{
		"Key1": NewString("Value1"),
		"Key2": NewList(NewString("Value2"), NewString("Value3")),
	}))
	require.Equal(t, true, b.IsMap())
	require.EqualValues(t, map[string]*TaToken{
		"Key1": NewString("Value1"),
		"Key2": NewList(NewString("Value2"), NewString("Value3")),
	}, b.Map())
}

func TestMapItem(t *testing.T) {
	tkn := NewMap(map[string]*TaToken{
		"Key1": NewBool(true),
		"Key2": NewDecimalFromInt(1),
		"Key3": NewString("Hello"),
	})
	require.Equal(t, NewBool(true), tkn.MapItem("Key1"))

	require.Equal(t, NewNull(), tkn.MapItem("Key4"))
}

func TestSetMapItem(t *testing.T) {
	tkn := NewMap(map[string]*TaToken{
		"Key1": NewBool(true),
		"Key2": NewDecimalFromInt(1),
		"Key3": NewString("Hello"),
	})
	tkn.SetMapItem("Key1", NewString("World"))
	require.Equal(t, NewString("World"), tkn.MapItem("Key1"))

	tkn.SetMapItem("Key4", NewString("Foo"))
	require.Equal(t, NewString("Foo"), tkn.MapItem("Key4"))
}

func TestStringify(t *testing.T) {
	tkn := New("+", NewDecimalFromInt(1), NewDecimalFromInt(2))
	require.Equal(t, "(+ 1 2)", tkn.Stringify())

	tkn = &TaToken{
		String: "noop",
		Kind:   Token,
	}
	require.Equal(t, "(noop)", tkn.Stringify())

	tkn = New("+", New("-", NewDecimalFromInt(1), NewDecimalFromInt(2)), NewDecimalFromInt(3))
	require.Equal(t, "(+ (- 1 2) 3)", tkn.Stringify())

	tkn = New("", New("-", NewDecimalFromInt(1), NewDecimalFromInt(2)), NewDecimalFromInt(3))
	require.Equal(t, "((- 1 2) 3)", tkn.Stringify())

	tkn = NewList(NewDecimalFromInt(1), NewDecimalFromInt(2), NewDecimalFromInt(3))
	require.Equal(t, "[1, 2, 3]", tkn.Stringify())

	tkn = NewString("Hello")
	require.Equal(t, "Hello", tkn.Stringify())

	time, err := time.Parse(time.UnixDate, "Mon Jan 2 15:04:05 MST 2006")
	require.NoError(t, err)

	tkn = NewTime(time)
	require.Equal(t, "2006-01-02T15:04:05Z", tkn.Stringify())

	tkn = NewBool(true)
	require.Equal(t, "true", tkn.Stringify())

	tkn = NewDecimal(decimal.New(145, 1))
	require.Equal(t, "14.5", tkn.Stringify())

	tkn = &TaToken{
		Kind: Map,
		Keys: []string{"Key1", "Key2"},
		Children: []*TaToken{
			&TaToken{
				Kind: Map,
				Keys: []string{"SubKey1"},
				Children: []*TaToken{
					&TaToken{
						String:  "14.5",
						Kind:    Decimal,
						Decimal: decimal.New(145, 1),
					},
				},
			},
			&TaToken{
				Kind:   String,
				String: "Hello",
			},
		},
	}

	require.Equal(t, "{Key1:{SubKey1:14.5}, Key2:Hello}", tkn.Stringify())
}

func TestEqual(t *testing.T) {
	now := time.Now()
	tests := []struct {
		a        *TaToken
		b        *TaToken
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
			NewMap(map[string]*TaToken{
				"Key1": NewBool(true),
				"Key2": NewDecimalFromInt(1),
				"Key3": NewString("Hello"),
			}),
			NewMap(map[string]*TaToken{
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
			&TaToken{},
			&TaToken{},
			false,
		},
		{
			NewMap(map[string]*TaToken{
				"Key1": NewBool(true),
				"Key3": NewString("Hello"),
			}),
			NewMap(map[string]*TaToken{
				"Key2": NewDecimalFromInt(1),
				"Key1": NewBool(true),
				"Key3": NewString("Hello"),
			}),
			false,
		},
		{
			NewMap(map[string]*TaToken{
				"Key1": NewBool(true),
				"Key2": NewDecimalFromInt(2),
				"Key3": NewString("Hello"),
			}),
			NewMap(map[string]*TaToken{
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
	tkn := New("+", New("1"), New("2"))
	require.EqualValues(t, []Kind{Decimal, Decimal}, Arguments(tkn.Children))

	tkn = New("+", New("Hello"), New("1"))
	require.EqualValues(t, []Kind{String, Decimal}, Arguments(tkn.Children))
}

func TestToHumanReadable(t *testing.T) {
	tkn := New("+", New("1"), New("2"))
	require.Equal(t, "1, 2", TokenArguments(tkn.Children).ToHumanReadable())
}

func TestMarshaling(t *testing.T) {
	block1 := NewMap(map[string]*TaToken{
		"Key2": NewDecimalFromInt(1),
		"Key1": NewBool(true),
		"Key3": NewString("Hello"),
		"Key4": NewList(NewBool(false), NewMap(map[string]*TaToken{
			"SubKey1": NewDecimalFromInt(3),
			"SubKey2": NewTime(time.Now()),
		})),
	})
	b, err := json.Marshal(block1)
	require.NoError(t, err)

	var block2 TaToken

	require.NoError(t, json.Unmarshal(b, &block2))

	require.Equal(t, true, block1.Equal(&block2))
}

func TestSort(t *testing.T) {
	list := NewList(NewDecimalFromInt(4), NewDecimalFromInt(2), NewDecimalFromInt(3))
	sort.Sort(TokenArguments(list.Children))
	require.Equal(t, true, list.Equal(NewList(NewDecimalFromInt(2), NewDecimalFromInt(3), NewDecimalFromInt(4))))
}

func BenchmarkIsDecimal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		isDecimal("2007-01-02T00:00:00Z")
	}
}

func BenchmarkTimeParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Parse(time.RFC3339, "2007-01-02T00:00:00Z")
	}
}

func BenchmarkMarshaling(ba *testing.B) {
	block1 := NewMap(map[string]*TaToken{
		"Key2": NewDecimalFromInt(1),
		"Key1": NewBool(true),
		"Key3": NewString("Hello"),
		"Key4": NewList(NewBool(false), NewMap(map[string]*TaToken{
			"SubKey1": NewDecimalFromInt(3),
			"SubKey2": NewTime(time.Now()),
		})),
	})
	b, err := json.Marshal(block1)
	require.NoError(ba, err)
	for i := 0; i < ba.N; i++ {
		var block2 TaToken
		require.NoError(ba, json.Unmarshal(b, &block2))
		require.Equal(ba, true, block1.Equal(&block2))
	}
}
