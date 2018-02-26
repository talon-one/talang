package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
)

func mustFunc(result *block.Block, err error) string {
	if err != nil {
		panic(err)
	}
	return result.String
}

func getError(result interface{}, err error) error {
	return err
}

func TestAdd(t *testing.T) {
	require.Error(t, getError(Add.Func(nil)))
	require.Error(t, getError(Add.Func(nil, block.New("1"))))
	require.Equal(t, "3", mustFunc(Add.Func(nil, block.New("1"), block.New("2"))))
	require.Equal(t, "4.6", mustFunc(Add.Func(nil, block.New("1.2"), block.New("3.4"))))
	require.Equal(t, "-1.0", mustFunc(Add.Func(nil, block.New("1.2"), block.New("3.4"), block.New("-5.6"))))
}

func TestSub(t *testing.T) {
	require.Error(t, getError(Sub.Func(nil)))
	require.Error(t, getError(Sub.Func(nil, block.New("1"))))
	require.Equal(t, "-1", mustFunc(Sub.Func(nil, block.New("1"), block.New("2"))))
	require.Equal(t, "-2.2", mustFunc(Sub.Func(nil, block.New("1.2"), block.New("3.4"))))
	require.Equal(t, "3.4", mustFunc(Sub.Func(nil, block.New("1.2"), block.New("3.4"), block.New("-5.6"))))
}

func TestMul(t *testing.T) {
	require.Error(t, getError(Mul.Func(nil)))
	require.Error(t, getError(Mul.Func(nil, block.New("1"))))
	require.Equal(t, "2", mustFunc(Mul.Func(nil, block.New("1"), block.New("2"))))
	require.Equal(t, "4.08", mustFunc(Mul.Func(nil, block.New("1.2"), block.New("3.4"))))
	require.Equal(t, "-22.848", mustFunc(Mul.Func(nil, block.New("1.2"), block.New("3.4"), block.New("-5.6"))))
}

func TestDiv(t *testing.T) {
	require.Error(t, getError(Div.Func(nil)))
	require.Error(t, getError(Div.Func(nil, block.New("1"))))
	require.Equal(t, "0.5", mustFunc(Div.Func(nil, block.New("1"), block.New("2"))))
	require.Equal(t, "0.3529411764705882", mustFunc(Div.Func(nil, block.New("1.2"), block.New("3.4"))))
	require.Equal(t, "-0.06302521008403361", mustFunc(Div.Func(nil, block.New("1.2"), block.New("3.4"), block.New("-5.6"))))
}

func TestMod(t *testing.T) {
	require.Error(t, getError(Mod.Func(nil)))
	require.Error(t, getError(Mod.Func(nil, block.New("1"))))
	require.Equal(t, "0", mustFunc(Mod.Func(nil, block.New("2"), block.New("1"))))
	require.Equal(t, "1", mustFunc(Mod.Func(nil, block.New("3"), block.New("2"))))
	require.Equal(t, "1", mustFunc(Mod.Func(nil, block.New("4"), block.New("3"), block.New("2"))))
}

func TestFloor(t *testing.T) {
	require.Error(t, getError(Floor.Func(nil)))
	require.Equal(t, "2", mustFunc(Floor.Func(nil, block.New("2"))))
	require.Equal(t, "2", mustFunc(Floor.Func(nil, block.New("2.4"))))
	require.Equal(t, "2", mustFunc(Floor.Func(nil, block.New("2.5"))))
	require.Equal(t, "2", mustFunc(Floor.Func(nil, block.New("2.9"))))
	require.Equal(t, "-3", mustFunc(Floor.Func(nil, block.New("-2.7"))))
	require.Equal(t, "-2", mustFunc(Floor.Func(nil, block.New("-2"))))
}

func TestCeil(t *testing.T) {
	require.Error(t, getError(Ceil.Func(nil)))
	require.Equal(t, "2", mustFunc(Ceil.Func(nil, block.New("2"))))
	require.Equal(t, "3", mustFunc(Ceil.Func(nil, block.New("2.4"))))
	require.Equal(t, "3", mustFunc(Ceil.Func(nil, block.New("2.9"))))
	require.Equal(t, "-2", mustFunc(Ceil.Func(nil, block.New("-2.7"))))
	require.Equal(t, "-2", mustFunc(Ceil.Func(nil, block.New("-2"))))
}

func TestAllOperations(t *testing.T) {
	AllOperations()
}

func TestStdFloor(t *testing.T) {
	require.Equal(t, float64(2), math.Floor(2))
	require.Equal(t, float64(2), math.Floor(2.4))
	require.Equal(t, float64(2), math.Floor(2.5))
	require.Equal(t, float64(2), math.Floor(2.9))
	require.Equal(t, float64(-3), math.Floor(-2.7))
	require.Equal(t, float64(-2), math.Floor(-2))
}

func TestStdCeil(t *testing.T) {
	require.Equal(t, float64(2), math.Ceil(2))
	require.Equal(t, float64(3), math.Ceil(2.4))
	require.Equal(t, float64(3), math.Ceil(2.5))
	require.Equal(t, float64(3), math.Ceil(2.9))
	require.Equal(t, float64(-2), math.Ceil(-2.7))
	require.Equal(t, float64(-2), math.Ceil(-2))
}
