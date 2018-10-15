package interpreter_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/talon-one/talang/lexer"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/decimal"
	"github.com/talon-one/talang/interpreter"
	helpers "github.com/talon-one/talang/testhelpers"
	"github.com/talon-one/talang/token"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"(+ 1 1)", "2"},
		{"(* (+ 1 2) 3)", "9"},
		{"(/ (+ 1 2) 3)", "1"},
		{"(/ (- 6 1) 2)", "2.5"},
		{"(= 1 1)", "true"},
		{"1", "1"},
	}

	interp := helpers.MustNewInterpreterWithLogger()

	for _, test := range tests {
		require.Equal(t, test.expected, interp.MustLexAndEvaluate(test.input).String, "Error in test `%s'", test.input)
	}
}

func TestInterpreterInvalidTerm(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.Error(t, interp.Evaluate(nil))
	require.Error(t, interp.Evaluate(&token.TaToken{}))
}

func TestOverloading(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RemoveAllFunctions())
	require.IsType(t, interpreter.FunctionNotFoundError{}, getError(interp.LexAndEvaluate("(FN 1 2)")))
	require.IsType(t, interpreter.FunctionNotFoundError{}, getError(interp.LexAndEvaluate("(FN A B)")))

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "FN",
			IsVariadic: false,
			Arguments: []token.Kind{
				token.Decimal,
				token.Decimal,
			},
			Returns: token.Decimal,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewDecimal(decimal.Add(args[0].Decimal, args[1].Decimal)), nil
		},
	})
	require.Equal(t, token.NewDecimalFromInt(3).Decimal, interp.MustLexAndEvaluate("(FN 1 2)").Decimal)
	require.IsType(t, interpreter.FunctionNotFoundError{}, getError(interp.LexAndEvaluate("(FN A B)")))

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "FN",
			IsVariadic: false,
			Arguments: []token.Kind{
				token.String,
				token.String,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			var sb strings.Builder
			sb.WriteString(args[0].String)
			sb.WriteString(args[1].String)
			return token.NewString(sb.String()), nil
		},
	})
	require.Equal(t, token.NewDecimalFromInt(3), interp.MustLexAndEvaluate("(FN 1 2)"))
	require.Equal(t, token.NewString("AB"), interp.MustLexAndEvaluate("(FN A B)"))
}

func TestOverloadingNested(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "fn1",
			IsVariadic: false,
			Arguments: []token.Kind{
				token.Decimal,
				token.Decimal,
			},
			Returns: token.Decimal,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewDecimal(decimal.Add(args[0].Decimal, args[1].Decimal)), nil
		},
	})

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "fn1",
			IsVariadic: false,
			Arguments: []token.Kind{
				token.String,
				token.String,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			var sb strings.Builder
			sb.WriteString(args[0].String)
			sb.WriteString(args[1].String)
			return token.NewString(sb.String()), nil
		},
	})

	nestedFuncCounter := 0
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "fn2",
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			nestedFuncCounter++
			return token.NewString(fmt.Sprintf("%d", nestedFuncCounter)), nil
		},
	})

	require.Equal(t, token.NewString("1C"), interp.MustLexAndEvaluate("(fn1 (fn2) C)"))
	require.Equal(t, 1, nestedFuncCounter)
}

func TestLists(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	result := interp.MustLexAndEvaluate("(list 1 2 3)")
	require.Equal(t, true, result.IsList())
	require.Equal(t, "", result.String)
	require.Equal(t, 3, len(result.Children))
}

func TestDoubleFuncCall(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	fn1Runned := false
	fn2Runned := false
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn1",
			Arguments: []token.Kind{
				token.Atom,
				token.Atom,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			fn1Runned = true
			return token.NewString("A"), nil
		},
	})
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn2",
			Arguments: []token.Kind{
				token.Atom,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			fn2Runned = true
			return token.NewString("B"), nil
		},
	})
	interp.MustLexAndEvaluate("(fn1 fn2 1)")
	require.Equal(t, true, fn1Runned)
	require.Equal(t, false, fn2Runned)
}

type customDataType struct {
	s string
}

func (d *customDataType) MarshalTaToken() (*token.TaToken, error) {
	return token.NewString(d.s), nil
}

func (d *customDataType) UnmarshalTaToken(tkn *token.TaToken) error {
	d.s = tkn.String
	return nil
}

func TestGenericSet(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	tests := []struct {
		input    interface{}
		expected *token.TaToken
	}{
		{"String", token.NewString("String")},
		{false, token.NewBool(false)},
		{123, token.NewDecimal(decimal.NewFromInt(123))},
	}

	for _, test := range tests {
		require.NoError(t, interp.GenericSet("Key", test.input), "Failed for %v", test.input)
		require.Equal(t, test.expected, interp.MustLexAndEvaluate("(. Key)"), "Failed for %v", test.input)
	}

	// struct
	require.NoError(t, interp.GenericSet("Key", struct {
		Str1 string
		Int2 int
	}{
		Str1: "Test",
		Int2: 1,
	}))
	require.Equal(t, "Test", interp.MustLexAndEvaluate("(. Key Str1)").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate("(. Key Int2)").String)

	// struct with a struct
	require.NoError(t, interp.GenericSet("Key", struct {
		Str1    string
		Struct2 struct {
			Str3 string
			Int4 int
		}
	}{
		Str1: "Test",
		Struct2: struct {
			Str3 string
			Int4 int
		}{
			Str3: "Hello",
			Int4: 1,
		},
	}))
	require.Equal(t, "Test", interp.MustLexAndEvaluate(". Key Str1").String)
	require.Equal(t, "Hello", interp.MustLexAndEvaluate(". Key Struct2 Str3").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate(". Key Struct2 Int4").String)

	// struct ptr
	require.NoError(t, interp.GenericSet("Key", &struct {
		Str1 string
		Int2 int
	}{
		Str1: "Test",
		Int2: 1,
	}))
	require.Equal(t, "Test", interp.MustLexAndEvaluate("(. Key Str1)").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate("(. Key Int2)").String)

	// map
	require.NoError(t, interp.GenericSet("Key", map[string]interface{}{
		"Str1": "Test",
		"Int2": 1,
	}))
	require.Equal(t, "Test", interp.MustLexAndEvaluate("(. Key Str1)").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate("(. Key Int2)").String)

	require.Error(t, interp.GenericSet("Key", map[int]interface{}{
		1: "Test",
		2: 1,
	}))

	// slice
	require.NoError(t, interp.GenericSet("Key", []interface{}{"Hello", true}))
	require.Equal(t, "Hello", interp.MustLexAndEvaluate("item (. Key) 0").String)
	require.Equal(t, true, interp.MustLexAndEvaluate("item (. Key) 1").Bool)

	// decimal
	require.NoError(t, interp.GenericSet("Key", decimal.NewFromInt(1)))
	require.Equal(t, "1", interp.MustLexAndEvaluate("(. Key)").String)
}

func TestGenericGet(t *testing.T) {
	t.Run("Primitives", func(t *testing.T) {
		type dataStruct struct {
			Str     string
			Int     int
			Int8    int8
			Int16   int16
			Int32   int32
			Int64   int64
			UInt    uint
			UInt8   uint8
			UInt16  uint16
			UInt32  uint32
			UInt64  uint64
			Float32 float32
			Float64 float64
			Bool    bool
		}

		interp := helpers.MustNewInterpreterWithLogger()
		interp.MustLexAndEvaluate(`(set Str Hello)`)
		interp.MustLexAndEvaluate(`(set Int 1)`)
		interp.MustLexAndEvaluate(`(set Int8 1)`)
		interp.MustLexAndEvaluate(`(set Int16 1)`)
		interp.MustLexAndEvaluate(`(set Int32 1)`)
		interp.MustLexAndEvaluate(`(set Int64 1)`)
		interp.MustLexAndEvaluate(`(set UInt 1)`)
		interp.MustLexAndEvaluate(`(set UInt8 1)`)
		interp.MustLexAndEvaluate(`(set UInt16 1)`)
		interp.MustLexAndEvaluate(`(set UInt32 1)`)
		interp.MustLexAndEvaluate(`(set UInt64 1)`)
		interp.MustLexAndEvaluate(`(set Float32 1.2)`)
		interp.MustLexAndEvaluate(`(set Float64 1.2)`)
		interp.MustLexAndEvaluate(`(set Bool true)`)

		t.Run("Combined", func(t *testing.T) {
			var data dataStruct
			require.NoError(t, interp.GenericGet("", &data))

			require.Equal(t, "Hello", data.Str)
			require.Equal(t, int(1), data.Int)
			require.Equal(t, int8(1), data.Int8)
			require.Equal(t, int16(1), data.Int16)
			require.Equal(t, int32(1), data.Int32)
			require.Equal(t, int64(1), data.Int64)
			require.Equal(t, uint(1), data.UInt)
			require.Equal(t, uint8(1), data.UInt8)
			require.Equal(t, uint16(1), data.UInt16)
			require.Equal(t, uint32(1), data.UInt32)
			require.Equal(t, uint64(1), data.UInt64)
			require.Equal(t, float32(1.2), data.Float32)
			require.Equal(t, float64(1.2), data.Float64)
			require.Equal(t, true, data.Bool)
		})

		t.Run("Single", func(t *testing.T) {
			var data dataStruct
			require.NoError(t, interp.GenericGet("Str", &data.Str))
			require.Equal(t, "Hello", data.Str)
			require.NoError(t, interp.GenericGet("Int", &data.Int))
			require.Equal(t, int(1), data.Int)
			require.NoError(t, interp.GenericGet("Int8", &data.Int8))
			require.Equal(t, int8(1), data.Int8)
			require.NoError(t, interp.GenericGet("Int16", &data.Int16))
			require.Equal(t, int16(1), data.Int16)
			require.NoError(t, interp.GenericGet("Int32", &data.Int32))
			require.Equal(t, int32(1), data.Int32)
			require.NoError(t, interp.GenericGet("Int64", &data.Int64))
			require.Equal(t, int64(1), data.Int64)
			require.NoError(t, interp.GenericGet("UInt", &data.UInt))
			require.Equal(t, uint(1), data.UInt)
			require.NoError(t, interp.GenericGet("UInt8", &data.UInt8))
			require.Equal(t, uint8(1), data.UInt8)
			require.NoError(t, interp.GenericGet("UInt16", &data.UInt16))
			require.Equal(t, uint16(1), data.UInt16)
			require.NoError(t, interp.GenericGet("UInt32", &data.UInt32))
			require.Equal(t, uint32(1), data.UInt32)
			require.NoError(t, interp.GenericGet("UInt64", &data.UInt64))
			require.Equal(t, uint64(1), data.UInt64)
			require.NoError(t, interp.GenericGet("Float32", &data.Float32))
			require.Equal(t, float32(1.2), data.Float32)
			require.NoError(t, interp.GenericGet("Float64", &data.Float64))
			require.Equal(t, float64(1.2), data.Float64)
			require.NoError(t, interp.GenericGet("Bool", &data.Bool))
			require.Equal(t, true, data.Bool)
		})
	})

	t.Run("Pointer", func(t *testing.T) {
		var data struct {
			Str  *string
			Int  *int
			Bool *bool
		}
		interp := helpers.MustNewInterpreterWithLogger()
		interp.MustLexAndEvaluate(`(set Str Hello)`)
		interp.MustLexAndEvaluate(`(set Int 1)`)
		interp.MustLexAndEvaluate(`(set Bool true)`)

		require.NoError(t, interp.GenericGet("", &data))

		require.Equal(t, "Hello", *data.Str)
		require.Equal(t, int(1), *data.Int)
		require.Equal(t, true, *data.Bool)
	})

	t.Run("Slice", func(t *testing.T) {
		var data struct {
			Slice1 []string
		}
		interp := helpers.MustNewInterpreterWithLogger()
		interp.MustLexAndEvaluate(`(set Slice1 (list Hello World))`)

		require.NoError(t, interp.GenericGet("", &data))

		require.EqualValues(t, []string{"Hello", "World"}, data.Slice1)
	})

	t.Run("Map", func(t *testing.T) {
		var data struct {
			Map1 map[string]string
		}
		interp := helpers.MustNewInterpreterWithLogger()
		interp.MustLexAndEvaluate(`(set Map1 Key1 World)`)
		interp.MustLexAndEvaluate(`(set Map1 Key2 Universe)`)

		require.NoError(t, interp.GenericGet("", &data))

		require.Equal(t, "World", data.Map1["Key1"])
		require.Equal(t, "Universe", data.Map1["Key2"])
	})

	t.Run("Struct", func(t *testing.T) {
		type subStruct struct {
			Str1 string
			Int2 int
		}
		var data struct {
			Sub1 subStruct
		}
		interp := helpers.MustNewInterpreterWithLogger()
		interp.MustLexAndEvaluate(`(set Sub1 Str1 World)`)
		interp.MustLexAndEvaluate(`(set Sub1 Int2 1)`)

		require.NoError(t, interp.GenericGet("", &data))

		require.Equal(t, "World", data.Sub1.Str1)
		require.Equal(t, 1, data.Sub1.Int2)
	})

	t.Run("Decimal", func(t *testing.T) {
		var dec decimal.Decimal
		interp := helpers.MustNewInterpreterWithLogger()
		interp.MustLexAndEvaluate(`(set Decimal 1)`)

		require.NoError(t, interp.GenericGet("Decimal", &dec))

		require.Equal(t, "1", dec.String())
	})

	t.Run("Custom", func(t *testing.T) {
		var data struct {
			Custom customDataType
		}
		interp := helpers.MustNewInterpreterWithLogger()
		interp.MustLexAndEvaluate(`(set Custom 1)`)

		require.NoError(t, interp.GenericGet("", &data))

		require.Equal(t, "1", data.Custom.s)
	})

	t.Run("CustomPTR", func(t *testing.T) {
		var data struct {
			Custom *customDataType
		}
		interp := helpers.MustNewInterpreterWithLogger()
		interp.MustLexAndEvaluate(`(set Custom 1)`)

		require.NoError(t, interp.GenericGet("", &data))

		require.Equal(t, "1", data.Custom.s)
	})
}

func TestMustEvaluate(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	b := lexer.MustLex("(panic)")
	require.Panics(t, func() { interp.MustEvaluate(b) })
}
func TestMustLexAndEvaluate(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.Panics(t, func() { interp.MustLexAndEvaluate("(panic)") })
}

func TestEvaluateResultIsBlock(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`FN Hello World`,
			nil,
			token.New("FN", token.NewString("Hello"), token.NewString("World")),
		},
	)
}

func TestModifiesInput(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterFunction(
		interpreter.TaFunction{
			CommonSignature: interpreter.CommonSignature{
				Name: "fn",
				Arguments: []token.Kind{
					token.List,
				},
				Returns: token.Any,
			},
			Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
				args[0].Children[0] = token.NewDecimalFromInt(1000)
				return nil, nil
			},
		},
		interpreter.TaFunction{
			CommonSignature: interpreter.CommonSignature{
				Name: "fn",
				Arguments: []token.Kind{
					token.Decimal,
				},
				Returns: token.Any,
			},
			Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
				args[0] = token.NewDecimalFromInt(1000)
				return nil, nil
			},
		},
	))

	interp.Binding = token.NewMap(map[string]*token.TaToken{
		"List1": token.NewList(token.NewDecimalFromInt(0), token.NewDecimalFromInt(1)),
		"Int1":  token.NewDecimalFromInt(0),
	})

	interp.MustLexAndEvaluate("(fn (. List1))")
	interp.MustLexAndEvaluate("(fn (. Int1))")

	require.Equal(t, true, interp.Get("List1").Equal(token.NewList(token.NewDecimalFromInt(0), token.NewDecimalFromInt(1))))
	require.Equal(t, true, interp.Get("Int1").Equal(token.NewDecimalFromInt(0)))
}

func BenchmarkInterpreter(b *testing.B) {
	tests := []struct {
		input    string
		expected string
	}{
		{"(+ 1 1)", "2"},
		{"(* (+ 1 2) 3)", "9"},
		{"(/ (+ 1 2) 3)", "1"},
		{"(/ (- 6 1) 2)", "2.5"},
		{"(= 1 1)", "true"},
	}

	interp := helpers.MustNewInterpreter()

	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			require.Equal(b, test.expected, interp.MustLexAndEvaluate(test.input).String)
		}
	}
}

func TestDryRun(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.IsDryRun = true

	interp.MustRegisterFunction(
		interpreter.TaFunction{
			CommonSignature: interpreter.CommonSignature{
				Name: "fn",
				Arguments: []token.Kind{
					token.Atom,
				},
				Returns: token.Atom,
			},
			Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
				panic("Function should have not been run")
				return nil, nil
			},
		},
	)

	interp.MustRegisterTemplate(
		interpreter.TaTemplate{
			CommonSignature: interpreter.CommonSignature{
				Name: "tmpl",
				Arguments: []token.Kind{
					token.Atom,
				},
				Returns: token.Atom,
			},
			Template: *lexer.MustLex("(panic Template should have not been run)"),
		},
	)

	parsedToken := lexer.MustLex("(fn (+ 1 2))")
	interp.MustEvaluate(parsedToken)

	parsedToken = lexer.MustLex("(! tmpl (+ 1 2))")
	interp.MustEvaluate(parsedToken)
}

func TestMultipleFuncCall(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"((set Integer 2) (+ (. Integer) 1))",
			nil,
			token.NewDecimalFromInt(3),
		},
	)
}

func TestDeepAbort(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.MaxRecursiveLevel = new(int)
	*interp.MaxRecursiveLevel = 10

	_, err := interp.LexAndEvaluate("(+ 1 (+ 2 (+ 3 (+ 4 (+ 5 (+ 6 (+ 7 (+ 8 (+ 9 (+ 10 (+ 11 (+ 12 (+ 13 (+ 14 15))))))))))))))")
	require.Error(t, err)
	//require.EqualError(t, err, interpreter.MaxRecursiveLevelReachedError{*interp.MaxRecursiveLevel}.Error())

	// ErrorTree:
	// +--- FunctionNotFound (+ 1 (+2 ...))
	//      +--- FunctionNotRan (+ 1 (+ 2...))
	//           +--- FunctionNotFound (+ 1 (+2 ...))
	//                +--- FunctionNotRan (+ 1 (+ 2...))

	// funcNotFoundErr, ok := err.(interpreter.FunctionNotFoundError)
	// require.Equal(t, true, ok, "error is not a FunctionNotFoundError was %T", err)
	// require.Equal(t, 1, len(funcNotFoundErr.CollectedErrors), "No errors in FunctionNotFoundError")

	// funcNotRanErr, ok := funcNotFoundErr.CollectedErrors[0].(interpreter.FunctionNotRanError)
	// require.Equal(t, true, ok, "error is not a FunctionNotRanError was %T", funcNotFoundErr.CollectedErrors[0])

	// funcNotFoundErr, ok = funcNotRanErr.Reason.(interpreter.FunctionNotFoundError)
	// require.Equal(t, true, ok, "error is not a FunctionNotFoundError was %T", funcNotRanErr.Reason)
	// require.Equal(t, 1, len(funcNotFoundErr.CollectedErrors), "No errors in FunctionNotFoundError")

	// _, ok = funcNotFoundErr.CollectedErrors[0].(interpreter.FunctionError)
	// require.Equal(t, true, ok, "error is not a FunctionError was %T", funcNotFoundErr.CollectedErrors[0])
}

func TestTypeChecking(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	_, err := interp.LexAndEvaluate(`(+ "2" 2)`)
	require.Error(t, err)
	// require.Equal(t, fmt.Sprintf("Found no eval function for (+ \"2\" 2)\n  Expression (+ \"2\" 2) doesn't match '+(Decimal, Decimal, Decimal...)Decimal'\n  Expression (+ \"2\" 2) doesn't match '+(String, String, String...)String'\n"), err.Error())
}
