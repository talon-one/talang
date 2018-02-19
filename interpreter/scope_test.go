package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

func TestScopeBinding(t *testing.T) {
	interp := MustNewInterpreter()
	interp.Set("RootKey", shared.Binding{
		Value: block.NewString("Root"),
	})

	require.Equal(t, "Root", interp.MustLexAndEvaluate("(. RootKey)").Text)

	scope := interp.NewScope()
	scope.Set("ScopeKey", shared.Binding{
		Value: block.NewString("Scope"),
	})
	require.Equal(t, "Root", scope.MustLexAndEvaluate("(. RootKey)").Text)

	_, err := interp.LexAndEvaluate("(. ScopeKey)")
	require.Error(t, err)
	require.Equal(t, "Scope", scope.MustLexAndEvaluate("(. ScopeKey)").Text)
}

func TestScopeFunctions(t *testing.T) {
	interp := MustNewInterpreter()
	interp.RegisterFunction(shared.TaSignature{
		Name: "fn1",
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello"), nil
		},
	})

	require.Equal(t, "Hello", interp.MustLexAndEvaluate("fn1").Text)

	scope := interp.NewScope()
	scope.RegisterFunction(shared.TaSignature{
		Name: "fn2",
		Func: func(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Bye"), nil
		},
	})
	require.Equal(t, "Hello", scope.MustLexAndEvaluate("fn1").Text)

	require.Equal(t, "fn2", interp.MustLexAndEvaluate("fn2").Text)
	require.Equal(t, "Bye", scope.MustLexAndEvaluate("fn2").Text)
}
