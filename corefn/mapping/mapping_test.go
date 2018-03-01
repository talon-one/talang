package mapping_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
	helpers "github.com/talon-one/talang/testhelpers"
)

func init() {
	interpreter.RegisterCoreFunction(
		interpreter.TaFunction{
			CommonSignature: interpreter.CommonSignature{
				Name:    "panic",
				Returns: block.AnyKind,
			},
			Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
				return nil, errors.New("panic")
			},
		},
	)
}

func TestKV(t *testing.T) {
	// tests need to be performed here because the testhelper is unable to compare maps properly

	interp := helpers.MustNewInterpreterWithLogger()
	result := interp.MustLexAndEvaluate("kv (Key1 Value1) (Key2 Value2)")
	require.Equal(t, true, result.IsMap())
	require.EqualValues(t, map[string]*block.Block{
		"Key1": block.NewString("Value1"),
		"Key2": block.NewString("Value2"),
	}, result.Map())

	result = interp.MustLexAndEvaluate("kv (Key1 Value1) (Key2 (list Value2 Value3))")
	require.Equal(t, true, result.IsMap())
	require.EqualValues(t, map[string]*block.Block{
		"Key1": block.NewString("Value1"),
		"Key2": block.NewList(block.NewString("Value2"), block.NewString("Value3")),
	}, result.Map())

	// other tests can run here
	helpers.RunTests(t,
		helpers.Test{
			"kv (Key1 Value1) (Key2 Value2 Value3)",
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			"kv (Key1 (+ 1 2))",
			nil,
			block.NewMap(map[string]*block.Block{
				"Key1": block.NewDecimalFromInt(3),
			}),
		},
		helpers.Test{
			"kv (Key1 (kv (SubKey Value)))",
			nil,
			block.NewMap(map[string]*block.Block{
				"Key1": block.NewMap(map[string]*block.Block{
					"SubKey": block.NewString("Value"),
				}),
			}),
		},
		helpers.Test{
			"kv (Key1 (panic))",
			nil,
			&helpers.Error{},
		},
	)
}
