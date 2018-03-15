package interpreter_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/interpreter"
	helpers "github.com/talon-one/talang/testhelpers"
	"github.com/talon-one/talang/token"
)

func getError(result interface{}, err error) error {
	return err
}

func TestRegisterFunction(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RemoveAllFunctions())
	// register a function
	require.NoError(t, interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "MyFN",
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").String)

	// try to register an already registered function
	require.Error(t, interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "myfn",
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Hello Universe"), nil
		}}))
	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").String)

	// update the function
	require.NoError(t, interp.UpdateFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "MyFn",
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Hello Galaxy"), nil
		}}))
	require.Equal(t, "Hello Galaxy", interp.MustLexAndEvaluate("myfn").String)

	// delete the function
	require.NoError(t, interp.RemoveFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "MyFN",
		},
	}))
	require.Equal(t, "myfn", interp.MustLexAndEvaluate("myfn").String)
}

func TestVariadicFunctionWith0Parameters(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "MyFN1",
			IsVariadic: true,
			Arguments: []token.Kind{
				token.Any,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn1").String)

	require.NoError(t, interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "MyFN2",
			IsVariadic: true,
			Returns:    token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn2").String)
}

func TestFuncWithWrongParameter(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "MyFN1",
			Arguments: []token.Kind{
				token.Any,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "myfn1", interp.MustLexAndEvaluate("myfn1").String)
}

func TestBinding(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	interp.Binding = token.NewMap(map[string]*token.TaToken{
		"Root1": token.NewMap(map[string]*token.TaToken{
			"Decimal": token.NewDecimalFromInt(2),
			"String":  token.NewString("Hello"),
			"List":    token.NewList(token.NewString("Item1"), token.NewString("Item2")),
			"Map": token.NewMap(map[string]*token.TaToken{
				"Decimal": token.NewDecimalFromInt(2),
				"String":  token.NewString("Hello"),
			}),
		}),
		"Root2": token.NewMap(map[string]*token.TaToken{}),
	})

	b := interp.MustLexAndEvaluate("(+ (. Root1 Decimal) 2)")
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "4", b.String)

	b = interp.MustLexAndEvaluate("(. Root1 String)")
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.String)

	b = interp.MustLexAndEvaluate("(. Root1 List)")
	require.Equal(t, true, b.IsList())

	b = interp.MustLexAndEvaluate("(. Root1 Map)")
	require.Equal(t, true, b.IsMap())

	b = interp.MustLexAndEvaluate("(. Root1 Map String)")
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.String)

	require.Error(t, getError(interp.LexAndEvaluate("(. Root1 Unknown)")))

	b = interp.MustLexAndEvaluate("(. Root2)")
	require.Equal(t, true, b.IsMap())

	require.Error(t, getError(interp.LexAndEvaluate("(. Root2 Decimal)")))
	require.Error(t, getError(interp.LexAndEvaluate("(. Root3)")))
	require.Error(t, getError(interp.LexAndEvaluate("(. Root3 Decimal)")))
}

func TestFuncInBinding(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	interp.Binding = token.NewMap(map[string]*token.TaToken{
		"Root": token.NewMap(map[string]*token.TaToken{
			"2": token.NewDecimalFromInt(2),
		}),
	})

	b := interp.MustLexAndEvaluate("(+ (. Root (+ 1 1)) 2)")
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "4", b.String)
}

func TestSetBinding(t *testing.T) {
	t.Run("RootLevel", func(t *testing.T) {
		interp := helpers.MustNewInterpreterWithLogger()
		interp.Binding = token.NewMap(map[string]*token.TaToken{
			"Root": token.NewMap(map[string]*token.TaToken{
				"Key": token.NewDecimalFromInt(1),
			}),
		})
		require.Equal(t, "1", interp.MustLexAndEvaluate(". Root Key").String)
		interp.MustLexAndEvaluate("set Root 2")
		require.Equal(t, "2", interp.MustLexAndEvaluate(". Root").String)
	})
	t.Run("DeepLevel", func(t *testing.T) {
		interp := helpers.MustNewInterpreterWithLogger()
		interp.Binding = token.NewMap(map[string]*token.TaToken{
			"Root": token.NewMap(map[string]*token.TaToken{
				"Key": token.NewDecimalFromInt(1),
			}),
		})
		interp.MustLexAndEvaluate("set Root Key 2")
		require.Equal(t, "2", interp.MustLexAndEvaluate(". Root Key").String)
	})
	t.Run("NotExistingRootLevel", func(t *testing.T) {
		interp := helpers.MustNewInterpreterWithLogger()
		require.Error(t, helpers.MustError(interp.LexAndEvaluate(". Root Key")))
		interp.MustLexAndEvaluate("set Root (kv (Key Hello))")
		require.Equal(t, "Hello", interp.MustLexAndEvaluate(". Root Key").String)
	})
	t.Run("NotExistingDeepLevel", func(t *testing.T) {
		interp := helpers.MustNewInterpreterWithLogger()
		require.Error(t, helpers.MustError(interp.LexAndEvaluate(". Root Key")))
		interp.MustLexAndEvaluate("set Root Key Hello")
		require.Equal(t, "Hello", interp.MustLexAndEvaluate(". Root Key").String)
	})
}

// Tests if a parent function can access the binding on a scoped interpreter
func TestRootFuncAccessScopeBinding(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn",
			Arguments: []token.Kind{
				token.String,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Hello " + args[0].String), nil
		},
	})

	scope := interp.NewScope()

	interp.Binding = token.NewMap(map[string]*token.TaToken{
		"Name": token.NewString("Joe"),
	})

	require.Equal(t, "Hello Joe", scope.MustLexAndEvaluate("fn (. Name)").String)
}

func TestFuncErrorInChild(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn1",
			Arguments: []token.Kind{
				token.String,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Test1"), nil
		},
	})

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn2",
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return nil, errors.New("SomeError")
		},
	})

	require.Error(t, getError(interp.LexAndEvaluate("fn1 (fn2)")))
}

func TestVariadicFunctionErrorInChild(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:       "fn1",
			IsVariadic: true,
			Arguments: []token.Kind{
				token.String,
			},
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewString("Test1"), nil
		},
	})

	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn2",
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return nil, errors.New("SomeError")
		},
	})

	require.Error(t, getError(interp.LexAndEvaluate("fn1 A (fn2)")))
}

// a function does not returns the correct type
func TestFunctionUnexpectedReturn(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "fn1",
			Returns: token.String,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return token.NewBool(true), nil
		},
	})

	require.Error(t, getError(interp.LexAndEvaluate("fn1")))
}

//  a function does not return a value and error
func TestFunctionNoReturnValue(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "fn1",
			Returns: token.Any,
		},
		Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
			return nil, nil
		},
	})

	require.Equal(t, token.Null, interp.MustLexAndEvaluate("fn1").Kind)
}
