package string

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
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

func TestContains(t *testing.T) {
	require.Error(t, getError(Contains.Func(nil, []*block.Block{})))
	require.Error(t, getError(Contains.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "false", mustFunc(Contains.Func(nil, []*block.Block{block.New("Hello World"), block.New("Universe")})))
	require.Equal(t, "false", mustFunc(Contains.Func(nil, []*block.Block{block.New("Hello World"), block.New("Hello"), block.New("Universe")})))
	require.Equal(t, "true", mustFunc(Contains.Func(nil, []*block.Block{block.New("Hello World"), block.New("Hello"), block.New("World")})))
}

func TestNotContains(t *testing.T) {
	require.Error(t, getError(NotContains.Func(nil, []*block.Block{})))
	require.Error(t, getError(NotContains.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(NotContains.Func(nil, []*block.Block{block.New("Hello World"), block.New("Universe")})))
	require.Equal(t, "false", mustFunc(NotContains.Func(nil, []*block.Block{block.New("Hello World"), block.New("Hello"), block.New("Universe")})))
	require.Equal(t, "false", mustFunc(NotContains.Func(nil, []*block.Block{block.New("Hello World"), block.New("Hello"), block.New("World")})))
}

func TestAllOperations(t *testing.T) {
	AllOperations()
}
