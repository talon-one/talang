package list

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

func TestList(t *testing.T) {
	require.Equal(t, &block.Block{
		Children: []*block.Block{},
		Kind:     block.BlockKind,
	}, mustFunc(List.Func(nil, []*block.Block{})))
	require.Equal(t, &block.Block{
		Children: []*block.Block{
			block.NewString("Hello"),
			block.NewString("World"),
		},
		Kind: block.BlockKind,
	}, mustFunc(List.Func(nil, []*block.Block{block.NewString("Hello"), block.NewString("World")})))
}

func TestHead(t *testing.T) {
	require.Error(t, getError(Head.Func(nil, []*block.Block{})))
	require.Equal(t, block.NewString("Hello"), mustFunc(Head.Func(nil, []*block.Block{block.NewString("Hello"), block.NewString("World")})))
}

func TestTail(t *testing.T) {
	require.Error(t, getError(Tail.Func(nil, []*block.Block{})))
	require.Equal(t, &block.Block{
		Children: []*block.Block{
			block.NewString("World"),
			block.NewString("and"),
			block.NewString("Universe"),
		},
		Kind: block.BlockKind,
	}, mustFunc(Tail.Func(nil, []*block.Block{block.NewString("Hello"), block.NewString("World"), block.NewString("and"), block.NewString("Universe")})))
}

func TestAllOperations(t *testing.T) {
	AllOperations()
}
