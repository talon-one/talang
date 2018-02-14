package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
)

func TestTemplate(t *testing.T) {
	interp := MustNewInterpreter()
	require.NoError(t, interp.SetTemplate("Template1", "(* 2 (. Variable1))"))

	var result *block.Block

	require.NoError(t, interp.GenericSet("Variable1", 1))
	result = interp.MustLexAndEvaluate("(+ 1 (! Template1))")
	require.Equal(t, true, result.IsDecimal())
	require.Equal(t, "3", result.Text)

	require.NoError(t, interp.GenericSet("Variable1", 2))
	result = interp.MustLexAndEvaluate("(+ 1 (! Template1))")
	require.Equal(t, true, result.IsDecimal())
	require.Equal(t, "5", result.Text)
}

func TestFormatedTemplate(t *testing.T) {
	interp := MustNewInterpreter()
	require.NoError(t, interp.SetTemplate("MultiplyWith2", "(* 2 (# 0))"))

	var result *block.Block

	result = interp.MustLexAndEvaluate("(+ 1 (! MultiplyWith2 2))")
	require.Equal(t, true, result.IsDecimal())
	require.Equal(t, "5", result.Text)

	result = interp.MustLexAndEvaluate("(+ 1 (! MultiplyWith2 4))")
	require.Equal(t, true, result.IsDecimal())
	require.Equal(t, "9", result.Text)
}
