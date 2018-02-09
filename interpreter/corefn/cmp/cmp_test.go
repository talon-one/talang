package cmp

import (
	"testing"
	"time"

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

func TestGreaterThanDecimal(t *testing.T) {
	require.Error(t, getError(GreaterThanDecimal.Func(nil, []*block.Block{})))
	require.Error(t, getError(GreaterThanDecimal.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(GreaterThanDecimal.Func(nil, []*block.Block{block.New("2"), block.New("1")})))
	require.Equal(t, "true", mustFunc(GreaterThanDecimal.Func(nil, []*block.Block{block.New("3"), block.New("1"), block.New("2")})))
	require.Equal(t, "false", mustFunc(GreaterThanDecimal.Func(nil, []*block.Block{block.New("0"), block.New("-1"), block.New("2")})))
	require.Equal(t, "false", mustFunc(GreaterThanDecimal.Func(nil, []*block.Block{block.New("1"), block.New("1")})))
}

func TestGreaterThanTime(t *testing.T) {
	now := time.Now()
	require.Error(t, getError(GreaterThanTime.Func(nil, []*block.Block{})))
	require.Error(t, getError(GreaterThanTime.Func(nil, []*block.Block{block.NewTime(now)})))
	require.Equal(t, "true", mustFunc(GreaterThanTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now.Add(-time.Second))})))
	require.Equal(t, "true", mustFunc(GreaterThanTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now.Add(-time.Second)), block.NewTime(now.Add(-time.Minute))})))
	require.Equal(t, "false", mustFunc(GreaterThanTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now.Add(-time.Second)), block.NewTime(now.Add(time.Minute))})))
	require.Equal(t, "false", mustFunc(GreaterThanTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now)})))
}

func TestLessThanDecimal(t *testing.T) {
	require.Error(t, getError(LessThanDecimal.Func(nil, []*block.Block{})))
	require.Error(t, getError(LessThanDecimal.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(LessThanDecimal.Func(nil, []*block.Block{block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(LessThanDecimal.Func(nil, []*block.Block{block.New("1"), block.New("2"), block.New("3")})))
	require.Equal(t, "false", mustFunc(LessThanDecimal.Func(nil, []*block.Block{block.New("1"), block.New("0"), block.New("2")})))
	require.Equal(t, "false", mustFunc(LessThanDecimal.Func(nil, []*block.Block{block.New("1"), block.New("1")})))
}

func TestLessThanTime(t *testing.T) {
	now := time.Now()
	require.Error(t, getError(LessThanTime.Func(nil, []*block.Block{})))
	require.Error(t, getError(LessThanTime.Func(nil, []*block.Block{block.NewTime(now)})))
	require.Equal(t, "true", mustFunc(LessThanTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now.Add(time.Second))})))
	require.Equal(t, "true", mustFunc(LessThanTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now.Add(time.Second)), block.NewTime(now.Add(time.Minute))})))
	require.Equal(t, "false", mustFunc(LessThanTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now.Add(-time.Second)), block.NewTime(now.Add(time.Second))})))
	require.Equal(t, "false", mustFunc(LessThanTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now)})))
}

func TestGreaterThanOrEqualDecimal(t *testing.T) {
	require.Error(t, getError(GreaterThanOrEqualDecimal.Func(nil, []*block.Block{})))
	require.Error(t, getError(GreaterThanOrEqualDecimal.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(GreaterThanOrEqualDecimal.Func(nil, []*block.Block{block.New("2"), block.New("1")})))
	require.Equal(t, "true", mustFunc(GreaterThanOrEqualDecimal.Func(nil, []*block.Block{block.New("3"), block.New("1"), block.New("2"), block.New("3")})))
	require.Equal(t, "false", mustFunc(GreaterThanOrEqualDecimal.Func(nil, []*block.Block{block.New("0"), block.New("-1"), block.New("2")})))
}

func TestGreaterThanOrEqualTime(t *testing.T) {
	now := time.Now()
	require.Error(t, getError(GreaterThanOrEqualTime.Func(nil, []*block.Block{})))
	require.Error(t, getError(GreaterThanOrEqualTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second))})))
	require.Equal(t, "true", mustFunc(GreaterThanOrEqualTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second * 2)), block.NewTime(now.Add(time.Second))})))
	require.Equal(t, "true", mustFunc(GreaterThanOrEqualTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second * 3)), block.NewTime(now.Add(time.Second)), block.NewTime(now.Add(time.Second * 2)), block.NewTime(now.Add(time.Second * 3))})))
	require.Equal(t, "false", mustFunc(GreaterThanOrEqualTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now.Add(-time.Second)), block.NewTime(now.Add(time.Second * 2))})))
}

func TestLessThanOrEqualDecimal(t *testing.T) {
	require.Error(t, getError(LessThanOrEqualDecimal.Func(nil, []*block.Block{})))
	require.Error(t, getError(LessThanOrEqualDecimal.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "true", mustFunc(LessThanOrEqualDecimal.Func(nil, []*block.Block{block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(LessThanOrEqualDecimal.Func(nil, []*block.Block{block.New("1"), block.New("2"), block.New("3"), block.New("1")})))
	require.Equal(t, "false", mustFunc(LessThanOrEqualDecimal.Func(nil, []*block.Block{block.New("1"), block.New("0"), block.New("2")})))
}

func TestLessThanOrEqualTime(t *testing.T) {
	now := time.Now()
	require.Error(t, getError(LessThanOrEqualTime.Func(nil, []*block.Block{})))
	require.Error(t, getError(LessThanOrEqualTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second))})))
	require.Equal(t, "true", mustFunc(LessThanOrEqualTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second)), block.NewTime(now.Add(time.Second * 2))})))
	require.Equal(t, "true", mustFunc(LessThanOrEqualTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second)), block.NewTime(now.Add(time.Second * 2)), block.NewTime(now.Add(time.Second * 3)), block.NewTime(now.Add(time.Second))})))
	require.Equal(t, "false", mustFunc(LessThanOrEqualTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second)), block.NewTime(now), block.NewTime(now.Add(time.Second * 2))})))
}

func TestBetweenDecimal(t *testing.T) {
	require.Error(t, getError(BetweenDecimal.Func(nil, []*block.Block{})))
	require.Error(t, getError(BetweenDecimal.Func(nil, []*block.Block{block.New("1")})))
	require.Error(t, getError(BetweenDecimal.Func(nil, []*block.Block{block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(BetweenDecimal.Func(nil, []*block.Block{block.New("1"), block.New("0"), block.New("2")})))
	require.Equal(t, "true", mustFunc(BetweenDecimal.Func(nil, []*block.Block{block.New("1"), block.New("1"), block.New("2")})))
	require.Equal(t, "true", mustFunc(BetweenDecimal.Func(nil, []*block.Block{block.New("1"), block.New("0"), block.New("1")})))
	require.Equal(t, "true", mustFunc(BetweenDecimal.Func(nil, []*block.Block{block.New("0"), block.New("1"), block.New("2"), block.New("0"), block.New("2")})))
	require.Equal(t, "false", mustFunc(BetweenDecimal.Func(nil, []*block.Block{block.New("3"), block.New("0"), block.New("2")})))
}

func TestBetweenTime(t *testing.T) {
	now := time.Now()
	require.Error(t, getError(BetweenTime.Func(nil, []*block.Block{})))
	require.Error(t, getError(BetweenTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second))})))
	require.Error(t, getError(BetweenTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second)), block.NewTime(now.Add(time.Second * 2))})))
	require.Equal(t, "true", mustFunc(BetweenTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second)), block.NewTime(now), block.NewTime(now.Add(time.Second * 2))})))
	require.Equal(t, "true", mustFunc(BetweenTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second)), block.NewTime(now.Add(time.Second)), block.NewTime(now.Add(time.Second * 2))})))
	require.Equal(t, "true", mustFunc(BetweenTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second)), block.NewTime(now), block.NewTime(now.Add(time.Second))})))
	require.Equal(t, "true", mustFunc(BetweenTime.Func(nil, []*block.Block{block.NewTime(now), block.NewTime(now.Add(time.Second)), block.NewTime(now.Add(time.Second * 2)), block.NewTime(now), block.NewTime(now.Add(time.Second * 2))})))
	require.Equal(t, "false", mustFunc(BetweenTime.Func(nil, []*block.Block{block.NewTime(now.Add(time.Second * 3)), block.NewTime(now), block.NewTime(now.Add(time.Second * 2))})))
}

func TestAllOperations(t *testing.T) {
	AllOperations()
}
