package interpreter_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/lexer"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestScopeBinding(t *testing.T) {
	// create an interpreter and set a binding
	interp := helpers.MustNewInterpreterWithLogger()
	interp.Set("RootKey", interpreter.Binding{
		Value: block.NewString("Root"),
	})

	// get the binding
	require.Equal(t, "Root", interp.MustLexAndEvaluate("(. RootKey)").String)

	// create a scope and set a binding ON the scope
	scope := interp.NewScope()
	scope.Set("ScopeKey", interpreter.Binding{
		Value: block.NewString("Scope"),
	})
	// check if the scope has the same binding as the root
	require.Equal(t, "Root", scope.MustLexAndEvaluate("(. RootKey)").String)

	// overwrite the binding on scope level
	scope.Set("RootKey", interpreter.Binding{
		Value: block.NewBool(true),
	})
	require.Equal(t, "true", scope.MustLexAndEvaluate("(. RootKey)").String)
	require.Equal(t, "Root", interp.MustLexAndEvaluate("(. RootKey)").String)

	_, err := interp.LexAndEvaluate("(. ScopeKey)")
	require.Error(t, err)
	require.Equal(t, "Scope", scope.MustLexAndEvaluate("(. ScopeKey)").String)
}

func TestScopeFunctions(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "fn1",
			Returns: block.StringKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Hello"), nil
		},
	})

	require.Equal(t, "Hello", interp.MustLexAndEvaluate("fn1").String)

	scope := interp.NewScope()
	scope.RegisterFunction(interpreter.TaFunction{
		CommonSignature: interpreter.CommonSignature{
			Name:    "fn2",
			Returns: block.StringKind,
		},
		Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
			return block.NewString("Bye"), nil
		},
	})
	require.Equal(t, "Hello", scope.MustLexAndEvaluate("fn1").String)

	require.Equal(t, "fn2", interp.MustLexAndEvaluate("fn2").String)
	require.Equal(t, "Bye", scope.MustLexAndEvaluate("fn2").String)
}

func TestScopeTemplates(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	require.NoError(t, interp.RegisterTemplate(interpreter.TaTemplate{
		CommonSignature: interpreter.CommonSignature{
			Name:    "Template1",
			Returns: block.StringKind,
		},
		Template: *lexer.MustLex("Hello"),
	}))
	require.Equal(t, "Hello", interp.MustLexAndEvaluate("! Template1").String)

	scope := interp.NewScope()

	require.NoError(t, scope.RegisterTemplate(interpreter.TaTemplate{
		CommonSignature: interpreter.CommonSignature{
			Name:    "Template2",
			Returns: block.StringKind,
		},
		Template: *lexer.MustLex("World"),
	}))
	require.Equal(t, "Hello World", scope.MustLexAndEvaluate(`+ (! Template1) " " (! Template2)`).String)

	require.Error(t, getError(interp.LexAndEvaluate("! Template2")))
}
