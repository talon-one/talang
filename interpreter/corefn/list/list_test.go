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
	require.Equal(t, block.New(""), mustFunc(List.Func(nil)))
	require.Equal(t, block.New("", block.NewString("Hello"), block.NewString("World")), mustFunc(List.Func(nil, block.NewString("Hello"), block.NewString("World"))))
}

func TestHead(t *testing.T) {
	require.Error(t, getError(Head.Func(nil)))
	require.Equal(t, block.NewString("Hello"), mustFunc(Head.Func(nil, block.NewString("Hello"), block.NewString("World"))))
}

func TestTail(t *testing.T) {
	require.Error(t, getError(Tail.Func(nil)))
	require.Equal(t, block.New("", block.NewString("World"), block.NewString("and"), block.NewString("Universe")), mustFunc(Tail.Func(nil, block.NewString("Hello"), block.NewString("World"), block.NewString("and"), block.NewString("Universe"))))
}

func TestDrop(t *testing.T) {
	require.Error(t, getError(Drop.Func(nil)))
	require.Equal(t, block.New("", block.NewString("Hello"), block.NewString("World"), block.NewString("and")), mustFunc(Drop.Func(nil, block.NewString("Hello"), block.NewString("World"), block.NewString("and"), block.NewString("Universe"))))
}

func TestItem(t *testing.T) {
	require.Error(t, getError(Item.Func(nil)))
	require.Error(t, getError(Item.Func(nil, block.New("", block.NewString("Hello"), block.NewString("World")), block.New("A"))))
	require.Error(t, getError(Item.Func(nil, block.New("", block.NewString("Hello"), block.NewString("World")), block.New("-1"))))
	require.Equal(t, block.NewString("Hello"), mustFunc(Item.Func(nil, block.New("", block.NewString("Hello"), block.NewString("World")), block.New("0"))))
	require.Equal(t, block.NewString("World"), mustFunc(Item.Func(nil, block.New("", block.NewString("Hello"), block.NewString("World")), block.New("1"))))
}

func TestAllOperations(t *testing.T) {
	AllOperations()
}
