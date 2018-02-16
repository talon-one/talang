package template

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
	"github.com/talon-one/talang/lexer"
)

func mustFunc(result *block.Block, err error) string {
	if err != nil {
		panic(err)
	}
	return result.String()
}

func getError(result interface{}, err error) error {
	return err
}

func TestTemplateBasic(t *testing.T) {
	var interp shared.Interpreter

	require.NoError(t, getError(SetTemplate.Func(&interp, block.NewString("Template1"), block.NewString("Hello World"))))
	require.NoError(t, getError(SetTemplate.Func(&interp, block.NewString("Template2"), block.NewString("Hello Universe"))))
	require.Equal(t, "Hello World", mustFunc(GetTemplate.Func(&interp, block.NewString("Template1"))))
	require.Equal(t, "Hello Universe", mustFunc(GetTemplate.Func(&interp, block.NewString("Template2"))))
	require.Error(t, getError(GetTemplate.Func(&interp, block.NewString("Template3"))))
	require.Error(t, getError(SetTemplate.Func(&interp, block.NewString("Template4"))))
}

func TestReplaceVariables(t *testing.T) {
	getCount := func(n int, err error) int {
		return n
	}
	require.Equal(t, 2, getCount(replaceVariables(lexer.MustLex("+ (# 1) (# 0)"), block.New("B"), block.New("A"))))
}

func TestReplaceVariablesInVariables(t *testing.T) {
	b := lexer.MustLex("+ (# 1) D")
	n, err := replaceVariables(b, lexer.MustLex("+ A B"), lexer.MustLex("+ (# 0) C"))
	if err != nil {
		panic(err)
	}

	require.Equal(t, 2, n)
	require.Equal(t, lexer.MustLex("+ (+ (+ A B) C) D"), b)
}

func TestAllOperations(t *testing.T) {
	AllOperations()
}
