package interpreter

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
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

	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)

	for _, test := range tests {
		require.Equal(t, test.expected, interp.MustLexAndEvaluate(test.input).Text, "Error in test `%s'", test.input)
	}
}

func TestInterpreterInvalidTerm(t *testing.T) {
	interp := MustNewInterpreter()
	require.Error(t, interp.Evaluate(nil))
	require.Error(t, interp.Evaluate(&block.Block{}))
}

func TestOverloading(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)
	require.NoError(t, interp.RemoveAllFunctions())
	require.Equal(t, "(+ 1 2)", interp.MustLexAndEvaluate("(+ 1 2)").String())
	require.Equal(t, "(+ A B)", interp.MustLexAndEvaluate("(+ A B)").String())

	interp.RegisterFunction(shared.TaSignature{
		Name:       "+",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewDecimal(args[0].Decimal.Add(args[0].Decimal, args[1].Decimal)), nil
		},
	})
	require.Equal(t, "3", interp.MustLexAndEvaluate("(+ 1 2)").String())
	require.Equal(t, "(+ A B)", interp.MustLexAndEvaluate("(+ A B)").String())

	interp.RegisterFunction(shared.TaSignature{
		Name:       "+",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.StringKind,
			block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.New(args[0].String() + args[1].String()), nil
		},
	})
	require.Equal(t, "3", interp.MustLexAndEvaluate("(+ 1 2)").String())
	require.Equal(t, "AB", interp.MustLexAndEvaluate("(+ A B)").String())
}

func TestOverloadingNested(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)

	interp.RegisterFunction(shared.TaSignature{
		Name:       "fn1",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.DecimalKind,
			block.DecimalKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewDecimal(args[0].Decimal.Add(args[0].Decimal, args[1].Decimal)), nil
		},
	})

	interp.RegisterFunction(shared.TaSignature{
		Name:       "fn1",
		IsVariadic: false,
		Arguments: []block.Kind{
			block.StringKind,
			block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.New(args[0].String() + args[1].String()), nil
		},
	})

	nestedFuncCounter := 0
	interp.RegisterFunction(shared.TaSignature{
		Name:       "fn2",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AnyKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			nestedFuncCounter++
			return block.NewString(fmt.Sprintf("%d", nestedFuncCounter)), nil
		},
	})

	require.Equal(t, "2C", interp.MustLexAndEvaluate("(fn1 (fn2) C)").String())
	require.Equal(t, 2, nestedFuncCounter)
}

func TestLists(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)
	result := interp.MustLexAndEvaluate("list 1 2 3")
	require.Equal(t, true, result.IsBlock())
	require.Equal(t, "", result.Text)
	require.Equal(t, 3, len(result.Children))
}

func TestDoubleFuncCall(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)

	fn1Runned := false
	fn2Runned := false
	interp.RegisterFunction(shared.TaSignature{
		Name: "fn1",
		Arguments: []block.Kind{
			block.AtomKind,
			block.AtomKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			fn1Runned = true
			return block.NewString("A"), nil
		},
	})
	interp.RegisterFunction(shared.TaSignature{
		Name: "fn2",
		Arguments: []block.Kind{
			block.AtomKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			fn2Runned = true
			return block.NewString("B"), nil
		},
	})
	interp.MustLexAndEvaluate("fn1 fn2 1")
	require.Equal(t, true, fn1Runned)
	require.Equal(t, false, fn2Runned)
}

func TestGenericSet(t *testing.T) {
	interp := MustNewInterpreter()

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
	require.Equal(t, "Test", interp.MustLexAndEvaluate(". Key Str1").Text)
	require.Equal(t, "1", interp.MustLexAndEvaluate(". Key Int2").Text)

	st := struct {
		Str1 string
		Int2 int
	}{
		Str1: "Test",
		Int2: 1,
	}
	require.NoError(t, interp.GenericSet("Key", &st))
	require.Equal(t, "Test", interp.MustLexAndEvaluate(". Key Str1").Text)
	require.Equal(t, "1", interp.MustLexAndEvaluate(". Key Int2").Text)
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

	interp := MustNewInterpreter()

	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			require.Equal(b, test.expected, interp.MustLexAndEvaluate(test.input).Text)
		}
	}
}
