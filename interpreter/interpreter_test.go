package interpreter_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/talon-one/talang/lexer"

	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"+ 1 1", "2"},
		{"* (+ 1 2) 3", "9"},
		{"/ (+ 1 2) 3", "1"},
		{"(/ (- 6 1) 2)", "2.5"},
		{"= 1 1", "true"},
	}

	interp := helpers.MustNewInterpreterWithLogger()

	for _, test := range tests {
		require.Equal(t, test.expected, interp.MustLexAndEvaluate(test.input).String, "Error in test `%s'", test.input)
	}
}

func TestInterpreterInvalidTerm(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.Error(t, interp.Evaluate(nil))
	require.Error(t, interp.Evaluate(&block.Block{}))
}

func TestOverloading(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RemoveAllFunctions())
	require.Equal(t, "(FN 1 2)", interp.MustLexAndEvaluate("(FN 1 2)").Stringify())
	require.Equal(t, "(FN A B)", interp.MustLexAndEvaluate("(FN A B)").Stringify())

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "FN",
			IsVariadic: false,
			Arguments: []block.Kind{
				block.DecimalKind,
				block.DecimalKind,
			},
			Returns: block.DecimalKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewDecimal(args[0].Decimal.Add(args[0].Decimal, args[1].Decimal)), nil
		},
	})
	require.Equal(t, "3", interp.MustLexAndEvaluate("(FN 1 2)").Stringify())
	require.Equal(t, "(FN A B)", interp.MustLexAndEvaluate("(FN A B)").Stringify())

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "FN",
			IsVariadic: false,
			Arguments: []block.Kind{
				block.StringKind,
				block.StringKind,
			},
			Returns: block.StringKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.New(args[0].Stringify() + args[1].Stringify()), nil
		},
	})
	require.Equal(t, "3", interp.MustLexAndEvaluate("(FN 1 2)").Stringify())
	require.Equal(t, "AB", interp.MustLexAndEvaluate("(FN A B)").Stringify())
}

func TestOverloadingNested(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "fn1",
			IsVariadic: false,
			Arguments: []block.Kind{
				block.DecimalKind,
				block.DecimalKind,
			},
			Returns: block.DecimalKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewDecimal(args[0].Decimal.Add(args[0].Decimal, args[1].Decimal)), nil
		},
	})

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "fn1",
			IsVariadic: false,
			Arguments: []block.Kind{
				block.StringKind,
				block.StringKind,
			},
			Returns: block.StringKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString(args[0].Stringify() + args[1].Stringify()), nil
		},
	})

	nestedFuncCounter := 0
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "fn2",
			Returns: block.StringKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			nestedFuncCounter++
			return block.NewString(fmt.Sprintf("%d", nestedFuncCounter)), nil
		},
	})

	require.Equal(t, "2C", interp.MustLexAndEvaluate("(fn1 (fn2) C)").Stringify())
	require.Equal(t, 2, nestedFuncCounter)
}

func TestLists(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	result := interp.MustLexAndEvaluate("list 1 2 3")
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
			Arguments: []block.Kind{
				block.AtomKind,
				block.AtomKind,
			},
			Returns: block.StringKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			fn1Runned = true
			return block.NewString("A"), nil
		},
	})
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn2",
			Arguments: []block.Kind{
				block.AtomKind,
			},
			Returns: block.StringKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			fn2Runned = true
			return block.NewString("B"), nil
		},
	})
	interp.MustLexAndEvaluate("fn1 fn2 1")
	require.Equal(t, true, fn1Runned)
	require.Equal(t, false, fn2Runned)
}

func TestGenericSet(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	tests := []struct {
		input    interface{}
		expected *block.Block
	}{
		{"String", block.NewString("String")},
		{false, block.NewBool(false)},
		{123, block.NewDecimal(decimal.New(123, 0))},
	}

	for _, test := range tests {
		require.NoError(t, interp.GenericSet("Key", test.input), "Failed for %v", test.input)
		require.Equal(t, test.expected, interp.MustLexAndEvaluate(". Key"), "Failed for %v", test.input)
	}

	require.NoError(t, interp.GenericSet("Key", struct {
		Str1 string
		Int2 int
	}{
		Str1: "Test",
		Int2: 1,
	}))
	require.Equal(t, "Test", interp.MustLexAndEvaluate(". Key Str1").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate(". Key Int2").String)

	st := struct {
		Str1 string
		Int2 int
	}{
		Str1: "Test",
		Int2: 1,
	}
	require.NoError(t, interp.GenericSet("Key", &st))
	require.Equal(t, "Test", interp.MustLexAndEvaluate(". Key Str1").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate(". Key Int2").String)

	require.NoError(t, interp.GenericSet("Key", map[string]interface{}{
		"Str1": "Test",
		"Int2": 1,
	}))
	require.Equal(t, "Test", interp.MustLexAndEvaluate(". Key Str1").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate(". Key Int2").String)

	require.Error(t, interp.GenericSet("Key", map[int]interface{}{
		1: "Test",
		2: 1,
	}))
}

func TestMustEvaluate(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterFunction(
		interpreter.TaFunction{
			CommonSignature: interpreter.CommonSignature{
				Name:    "panic",
				Returns: block.AnyKind,
			},
			Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
				return nil, errors.New("panic")
			},
		},
	))
	b := lexer.MustLex("panic")
	require.Panics(t, func() { interp.MustEvaluate(b) })
}
func TestMustLexAndEvaluate(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterFunction(
		interpreter.TaFunction{
			CommonSignature: interpreter.CommonSignature{
				Name:    "panic",
				Returns: block.AnyKind,
			},
			Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
				return nil, errors.New("panic")
			},
		},
	))
	require.Panics(t, func() { interp.MustLexAndEvaluate("panic") })
}

func TestEvaluateResultIsBlock(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`FN Hello World`,
			nil,
			block.New("FN", block.NewString("Hello"), block.NewString("World")),
		},
	)
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

	interp := helpers.MustNewInterpreterWithLogger()

	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			require.Equal(b, test.expected, interp.MustLexAndEvaluate(test.input).String)
		}
	}
}
