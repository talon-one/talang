package misc

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
)

func mustFunc(result *block.Block, err error) *block.Block {
	if err != nil {
		panic(err)
	}
	return result
}

func getError(result interface{}, err error) error {
	return err
}

func TestNoop(t *testing.T) {
	require.NoError(t, getError(Noop.Func(nil)))
}

func TestToString(t *testing.T) {
	require.Error(t, getError(ToString.Func(nil)))
	require.Equal(t, block.NewString("1"), mustFunc(ToString.Func(nil, block.New("1"))))
}
