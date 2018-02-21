package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

func TestScopeBinding(t *testing.T) {
	// create an interpreter and set a binding
	interp := MustNewInterpreter()
	interp.Set("RootKey", shared.Binding{
		Value: block.NewString("Root"),
	})

	// get the binding
	require.Equal(t, "Root", interp.MustLexAndEvaluate("(. RootKey)").Text)

	// create a scope and set a binding ON the scope
	scope := interp.NewScope()
	scope.Set("ScopeKey", shared.Binding{
		Value: block.NewString("Scope"),
	})
	// check if the scope has the same binding as the root
	require.Equal(t, "Root", scope.MustLexAndEvaluate("(. RootKey)").Text)

	// overwrite the binding on scope level
	scope.Set("RootKey", shared.Binding{
		Value: block.NewBool(true),
	})
	require.Equal(t, "true", scope.MustLexAndEvaluate("(. RootKey)").Text)
	require.Equal(t, "Root", interp.MustLexAndEvaluate("(. RootKey)").Text)

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

func TestScopeTemplates(t *testing.T) {
	interp := MustNewInterpreter()
	require.NoError(t, interp.SetTemplate("Template1", "Hello"))
	require.Equal(t, "Hello", interp.MustLexAndEvaluate("! Template1").Text)

	scope := interp.NewScope()

	require.NoError(t, scope.SetTemplate("Template2", "World"))
	require.Equal(t, "Hello World", scope.MustLexAndEvaluate(`+ (! Template1) " " (! Template2)`).Text)

	require.Error(t, getError(interp.LexAndEvaluate("! Template2")))
}
