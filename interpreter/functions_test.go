package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/interpreter/internal"
	"github.com/talon-one/talang/term"
)

func TestRegisterFunction(t *testing.T) {
	interp := MustNewInterpreter()

	// register a function
	require.NoError(t, interp.RegisterFunction("MyFN", func(interp *internal.Interpreter, args ...term.Term) (string, error) {
		return "Hello World", nil
	}))
	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").Text)

	// try to register an already registered function
	require.Error(t, interp.RegisterFunction("myfn", func(interp *internal.Interpreter, args ...term.Term) (string, error) {
		return "Hello Universe", nil
	}))
	require.Equal(t, "Hello World", interp.MustLexAndEvaluate("myfn").Text)

	// update the function
	require.NoError(t, interp.UpdateFunction("MyFn", func(interp *internal.Interpreter, args ...term.Term) (string, error) {
		return "Hello Galaxy", nil
	}))
	require.Equal(t, "Hello Galaxy", interp.MustLexAndEvaluate("myfn").Text)

	// delete the function
	require.NoError(t, interp.RemoveFunction("MyFn"))
	require.Equal(t, "myfn", interp.MustLexAndEvaluate("myfn").Text)

}
