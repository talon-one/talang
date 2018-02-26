package interpreter

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

func getError(result interface{}, err error) error {
	return err
}

func TestRegisterFunction(t *testing.T) {
	interp := mustNewInterpreterWithLogger()
	require.NoError(t, interp.RemoveAllFunctions())
	// register a function
	require.NoError(t, interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:    "MyFN",
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").String)

	// try to register an already registered function
	require.Error(t, interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:    "myfn",
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello Universe"), nil
		}}))
	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").String)

	// update the function
	require.NoError(t, interp.UpdateFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:    "MyFn",
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello Galaxy"), nil
		}}))
	require.Equal(t, "Hello Galaxy", interp.MustLexAndEvaluate("myfn").String)

	// delete the function
	require.NoError(t, interp.RemoveFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name: "MyFN",
		},
	}))
	require.Equal(t, "myfn", interp.MustLexAndEvaluate("myfn").String)
}

func TestVariadicFunctionWith0Parameters(t *testing.T) {
	interp := mustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:       "MyFN1",
			IsVariadic: true,
			Arguments: []block.Kind{
				block.AnyKind,
			},
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn1").String)

	require.NoError(t, interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:       "MyFN2",
			IsVariadic: true,
			Returns:    block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn2").String)
}

func TestFuncWithWrongParameter(t *testing.T) {
	interp := mustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name: "MyFN1",
			Arguments: []block.Kind{
				block.AnyKind,
			},
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "myfn1", interp.MustLexAndEvaluate("myfn1").String)
}

func TestBinding(t *testing.T) {
	interp := mustNewInterpreterWithLogger()

	interp.Binding["Root1"] = shared.Binding{
		Value: block.New("1"),
		Children: map[string]shared.Binding{
			"Decimal": shared.Binding{
				Value: block.New("2"),
			},
			"String": shared.Binding{
				Value: block.New("Hello"),
			},
			"List": shared.Binding{
				Value: block.New("", block.NewString("Item1"), block.NewString("Item2")),
			},
			"Map": shared.Binding{
				Children: map[string]shared.Binding{
					"Decimal": shared.Binding{
						Value: block.New("2"),
					},
					"String": shared.Binding{
						Value: block.New("Hello"),
					},
				},
			},
		},
	}
	interp.Binding["Root2"] = shared.Binding{}

	b := interp.MustLexAndEvaluate("(+ (. Root1) 2)")
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "3", b.String)
	b = interp.MustLexAndEvaluate("(+ (. Root1 Decimal) 2)")
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "4", b.String)

	b = interp.MustLexAndEvaluate("(. Root1 String)")
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.String)

	b = interp.MustLexAndEvaluate("(. Root1 List)")
	require.Equal(t, true, b.IsBlock())

	require.Error(t, getError(interp.LexAndEvaluate("(. Root1 Map)")))

	b = interp.MustLexAndEvaluate("(. Root1 Map String)")
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.String)

	require.Error(t, getError(interp.LexAndEvaluate("(. Root1 Unknown)")))

	require.Error(t, getError(interp.LexAndEvaluate("(. Root2)")))
	require.Error(t, getError(interp.LexAndEvaluate("(. Root2 Decimal)")))
	require.Error(t, getError(interp.LexAndEvaluate("(. Root3)")))
	require.Error(t, getError(interp.LexAndEvaluate("(. Root3 Decimal)")))
}

func TestFuncInBinding(t *testing.T) {
	interp := mustNewInterpreterWithLogger()

	interp.Binding["Root"] = shared.Binding{
		Children: map[string]shared.Binding{
			"2": shared.Binding{
				Value: block.New("2"),
			},
		},
	}

	b := interp.MustLexAndEvaluate("(+ (. Root (+ 1 1)) 2)")
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "4", b.String)
}

func TestFuncErrorInChild(t *testing.T) {
	interp := mustNewInterpreterWithLogger()
	interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name: "fn1",
			Arguments: []block.Kind{
				block.StringKind,
			},
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Test1"), nil
		},
	})

	interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name: "fn2",
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return nil, errors.New("SomeError")
		},
	})

	require.Error(t, getError(interp.LexAndEvaluate("fn1 (fn2)")))
}

func TestVariadicFunctionErrorInChild(t *testing.T) {
	interp := mustNewInterpreterWithLogger()
	interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:       "fn1",
			IsVariadic: true,
			Arguments: []block.Kind{
				block.StringKind,
			},
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Test1"), nil
		},
	})

	interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name: "fn2",
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return nil, errors.New("SomeError")
		},
	})

	require.Error(t, getError(interp.LexAndEvaluate("fn1 A (fn2)")))
}

// a function does not returns the correct type
func TestFunctionUnexpectedReturn(t *testing.T) {
	interp := mustNewInterpreterWithLogger()
	interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:    "fn1",
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewBool(true), nil
		},
	})

	require.Error(t, getError(interp.LexAndEvaluate("fn1")))
}

//  a function does not return a value and error
func TestFunctionNoReturnValue(t *testing.T) {
	interp := mustNewInterpreterWithLogger()
	interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:    "fn1",
			Returns: block.AnyKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return nil, nil
		},
	})

	require.Equal(t, block.NullKind, interp.MustLexAndEvaluate("fn1").Kind)
}
