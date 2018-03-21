package interpreter_test

import (
	"fmt"
	"testing"

	"github.com/talon-one/talang/lexer"

	"github.com/ericlagergren/decimal"
	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/interpreter"
	helpers "github.com/talon-one/talang/testhelpers"
	"github.com/talon-one/talang/token"
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
	require.Error(t, interp.Evaluate(&token.TaToken{}))
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
			Arguments: []token.Kind{
				token.Decimal,
				token.Decimal,
			},
			Returns: token.Decimal,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewDecimal(args[0].Decimal.Add(args[0].Decimal, args[1].Decimal)), nil
		},
	})
	require.Equal(t, "3", interp.MustLexAndEvaluate("(FN 1 2)").Stringify())
	require.Equal(t, "(FN A B)", interp.MustLexAndEvaluate("(FN A B)").Stringify())

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
			return token.New(args[0].Stringify() + args[1].Stringify()), nil
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
			Arguments: []token.Kind{
				token.Decimal,
				token.Decimal,
			},
			Returns: token.Decimal,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewDecimal(args[0].Decimal.Add(args[0].Decimal, args[1].Decimal)), nil
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
			return token.NewString(args[0].Stringify() + args[1].Stringify()), nil
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
	interp.MustLexAndEvaluate("fn1 fn2 1")
	require.Equal(t, true, fn1Runned)
	require.Equal(t, false, fn2Runned)
}

func TestGenericSet(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	tests := []struct {
		input    interface{}
		expected *token.TaToken
	}{
		{"String", token.NewString("String")},
		{false, token.NewBool(false)},
		{123, token.NewDecimal(decimal.New(123, 0))},
	}

	for _, test := range tests {
		require.NoError(t, interp.GenericSet("Key", test.input), "Failed for %v", test.input)
		require.Equal(t, test.expected, interp.MustLexAndEvaluate(". Key"), "Failed for %v", test.input)
	}

	// struct
	require.NoError(t, interp.GenericSet("Key", struct {
		Str1 string
		Int2 int
	}{
		Str1: "Test",
		Int2: 1,
	}))
	require.Equal(t, "Test", interp.MustLexAndEvaluate(". Key Str1").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate(". Key Int2").String)

	// struct ptr
	require.NoError(t, interp.GenericSet("Key", &struct {
		Str1 string
		Int2 int
	}{
		Str1: "Test",
		Int2: 1,
	}))
	require.Equal(t, "Test", interp.MustLexAndEvaluate(". Key Str1").String)
	require.Equal(t, "1", interp.MustLexAndEvaluate(". Key Int2").String)

	// map
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

	// slice
	require.NoError(t, interp.GenericSet("Key", []interface{}{"Hello", true}))
	require.Equal(t, "Hello", interp.MustLexAndEvaluate("item (. Key) 0").String)
	require.Equal(t, true, interp.MustLexAndEvaluate("item (. Key) 1").Bool)
}

func TestMustEvaluate(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	b := lexer.MustLex("panic")
	require.Panics(t, func() { interp.MustEvaluate(b) })
}
func TestMustLexAndEvaluate(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.Panics(t, func() { interp.MustLexAndEvaluate("panic") })
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

	interp.MustLexAndEvaluate("fn (. List1)")
	interp.MustLexAndEvaluate("fn (. Int1)")

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

func DisableTestDryRun(t *testing.T) {
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
	var evalToken token.TaToken
	token.Copy(&evalToken, parsedToken)
	interp.MustEvaluate(&evalToken)
	require.Equal(t, true, evalToken.Equal(parsedToken))

	parsedToken = lexer.MustLex("(tmpl (+ 1 2))")
	token.Copy(&evalToken, parsedToken)
	interp.MustEvaluate(&evalToken)
	require.Equal(t, true, evalToken.Equal(parsedToken))
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
