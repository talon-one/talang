package cmp

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

func TestEqual(t *testing.T) {
	require.Error(t, getError(Equal.Func(nil, []*block.Block{})))
	require.Error(t, getError(Equal.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "false", mustFunc(Equal.Func(nil, []*block.Block{block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(Equal.Func(nil, []*block.Block{block.New("1"), block.New("1"), block.New("1")})))
}

func TestNotEqual(t *testing.T) {
	require.Error(t, getError(NotEqual.Func(nil, []*block.Block{})))
	require.Error(t, getError(NotEqual.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(NotEqual.Func(nil, []*block.Block{block.New("1"), block.New("2")})))
	require.Equal(t, "false", mustFunc(NotEqual.Func(nil, []*block.Block{block.New("1"), block.New("1"), block.New("1")})))
}

func TestGreaterThan(t *testing.T) {
	require.Error(t, getError(GreaterThan.Func(nil, []*block.Block{})))
	require.Error(t, getError(GreaterThan.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(GreaterThan.Func(nil, []*block.Block{block.New("2"), block.New("1")})))
	require.Equal(t, "true", mustFunc(GreaterThan.Func(nil, []*block.Block{block.New("3"), block.New("1"), block.New("2")})))
	require.Equal(t, "false", mustFunc(GreaterThan.Func(nil, []*block.Block{block.New("0"), block.New("-1"), block.New("2")})))
	require.Equal(t, "false", mustFunc(GreaterThan.Func(nil, []*block.Block{block.New("1"), block.New("1")})))
}

func TestLessThan(t *testing.T) {
	require.Error(t, getError(LessThan.Func(nil, []*block.Block{})))
	require.Error(t, getError(LessThan.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(LessThan.Func(nil, []*block.Block{block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(LessThan.Func(nil, []*block.Block{block.New("1"), block.New("2"), block.New("3")})))
	require.Equal(t, "false", mustFunc(LessThan.Func(nil, []*block.Block{block.New("1"), block.New("0"), block.New("2")})))
	require.Equal(t, "false", mustFunc(LessThan.Func(nil, []*block.Block{block.New("1"), block.New("1")})))
}

func TestGreaterThanOrEqual(t *testing.T) {
	require.Error(t, getError(GreaterThanOrEqual.Func(nil, []*block.Block{})))
	require.Error(t, getError(GreaterThanOrEqual.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(GreaterThanOrEqual.Func(nil, []*block.Block{block.New("2"), block.New("1")})))
	require.Equal(t, "true", mustFunc(GreaterThanOrEqual.Func(nil, []*block.Block{block.New("3"), block.New("1"), block.New("2"), block.New("3")})))
	require.Equal(t, "false", mustFunc(GreaterThanOrEqual.Func(nil, []*block.Block{block.New("0"), block.New("-1"), block.New("2")})))
}

func TestLessThanOrEqual(t *testing.T) {
	require.Error(t, getError(LessThanOrEqual.Func(nil, []*block.Block{})))
	require.Error(t, getError(LessThanOrEqual.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(LessThanOrEqual.Func(nil, []*block.Block{block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(LessThanOrEqual.Func(nil, []*block.Block{block.New("1"), block.New("2"), block.New("3"), block.New("1")})))
	require.Equal(t, "false", mustFunc(LessThanOrEqual.Func(nil, []*block.Block{block.New("1"), block.New("0"), block.New("2")})))
}

func TestBetween(t *testing.T) {
	require.Error(t, getError(Between.Func(nil, []*block.Block{})))
	require.Error(t, getError(Between.Func(nil, []*block.Block{block.New("1")})))
	require.Error(t, getError(Between.Func(nil, []*block.Block{block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(Between.Func(nil, []*block.Block{block.New("1"), block.New("0"), block.New("2")})))
	require.Equal(t, "true", mustFunc(Between.Func(nil, []*block.Block{block.New("1"), block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(Between.Func(nil, []*block.Block{block.New("1"), block.New("0"), block.New("1")})))
	require.Equal(t, "true", mustFunc(Between.Func(nil, []*block.Block{block.New("0"), block.New("1"), block.New("2"), block.New("0"), block.New("2")})))
	require.Equal(t, "false", mustFunc(Between.Func(nil, []*block.Block{block.New("3"), block.New("0"), block.New("2")})))
}
