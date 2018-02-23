package interpreter

import (
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

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").Text)

	// try to register an already registered function
	require.Error(t, interp.RegisterFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:    "myfn",
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello Universe"), nil
		}}))
	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").Text)

	// update the function
	require.NoError(t, interp.UpdateFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name:    "MyFn",
			Returns: block.StringKind,
		},
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello Galaxy"), nil
		}}))
	require.Equal(t, "Hello Galaxy", interp.MustLexAndEvaluate("myfn").Text)

	// delete the function
	require.NoError(t, interp.RemoveFunction(shared.TaFunction{
		CommonSignature: shared.CommonSignature{
			Name: "MyFN",
		},
	}))
	require.Equal(t, "myfn", interp.MustLexAndEvaluate("myfn").Text)
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

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn1").Text)

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

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn2").Text)
}

func TestFunctionWithWrongParamter(t *testing.T) {
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

	require.Equal(t, "myfn1", interp.MustLexAndEvaluate("myfn1").Text)
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
	require.Equal(t, "3", b.Text)
	b = interp.MustLexAndEvaluate("(+ (. Root1 Decimal) 2)")
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "4", b.Text)

	b = interp.MustLexAndEvaluate("(. Root1 String)")
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.Text)

	b = interp.MustLexAndEvaluate("(. Root1 List)")
	require.Equal(t, true, b.IsBlock())

	require.Error(t, getError(interp.LexAndEvaluate("(. Root1 Map)")))

	b = interp.MustLexAndEvaluate("(. Root1 Map String)")
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.Text)

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
	require.Equal(t, "4", b.Text)
}
