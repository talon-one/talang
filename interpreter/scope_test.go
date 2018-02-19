package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
)

func TestScope(t *testing.T) {
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
