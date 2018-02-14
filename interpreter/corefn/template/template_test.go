package template

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/shared"
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
}
