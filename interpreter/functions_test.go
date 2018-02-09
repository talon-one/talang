package interpreter

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

func getError(result interface{}, err error) error {
	return err
}

func TestRegisterFunction(t *testing.T) {
	interp := MustNewInterpreter()
	require.NoError(t, interp.RemoveAllFunctions())
	// register a function
	require.NoError(t, interp.RegisterFunction(shared.TaSignature{
		Name: "MyFN",
		Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
			return block.NewString("Hello World"), nil
		},
	}))

	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").Text)

	// try to register an already registered function
	require.Error(t, interp.RegisterFunction(shared.TaSignature{
		Name: "myfn",
		Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
			return block.NewString("Hello Universe"), nil
		}}))
	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").Text)

	// update the function
	require.NoError(t, interp.UpdateFunction(shared.TaSignature{
		Name: "MyFn",
		Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
			return block.NewString("Hello Galaxy"), nil
		}}))
	require.Equal(t, "Hello Galaxy", interp.MustLexAndEvaluate("myfn").Text)

	// delete the function
	require.NoError(t, interp.RemoveFunction(shared.TaSignature{
		Name: "MyFN",
	}))
	require.Equal(t, "myfn", interp.MustLexAndEvaluate("myfn").Text)

}

func TestBinding(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)

	interp.Binding["Root"] = shared.Binding{
		Value: block.New("1"),
		Children: map[string]shared.Binding{
			"Child1": shared.Binding{
				Value: block.New("2"),
			},
			"Child2": shared.Binding{
				Value: block.New("Hello"),
			},
		},
	}

	b := interp.MustLexAndEvaluate("(+ (. Root) 2)")
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "3", b.Text)
	b = interp.MustLexAndEvaluate("(+ (. Root Child1) 2)")
	require.Equal(t, true, b.IsDecimal())
	require.Equal(t, "4", b.Text)

	b = interp.MustLexAndEvaluate("(. Root Child2)")
	require.Equal(t, true, b.IsString())
	require.Equal(t, "Hello", b.Text)

	require.Error(t, getError(interp.LexAndEvaluate("(. Root Child3)")))
}
