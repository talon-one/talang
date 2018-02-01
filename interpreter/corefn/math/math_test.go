package math

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/term"
)

func mustFunc(result string, err error) string {
	if err != nil {
		panic(err)
	}
	return result
}

func getError(result interface{}, err error) error {
	return err
}

func TestAdd(t *testing.T) {
	require.Error(t, getError(Add(nil)))
	require.Error(t, getError(Add(nil, term.New("1"))))
	require.Equal(t, "3", mustFunc(Add(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "4.6", mustFunc(Add(nil, term.New("1.2"), term.New("3.4"))))
	require.Equal(t, "-1", mustFunc(Add(nil, term.New("1.2"), term.New("3.4"), term.New("-5.6"))))
	require.Error(t, getError(Add(nil, term.New("1.2"), term.New("3.4"), term.New("A"))))
}

func TestSub(t *testing.T) {
	require.Error(t, getError(Sub(nil)))
	require.Error(t, getError(Sub(nil, term.New("1"))))
	require.Equal(t, "-1", mustFunc(Sub(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "-2.2", mustFunc(Sub(nil, term.New("1.2"), term.New("3.4"))))
	require.Equal(t, "3.4", mustFunc(Sub(nil, term.New("1.2"), term.New("3.4"), term.New("-5.6"))))
	require.Error(t, getError(Sub(nil, term.New("1.2"), term.New("3.4"), term.New("A"))))
}

func TestMul(t *testing.T) {
	require.Error(t, getError(Mul(nil)))
	require.Error(t, getError(Mul(nil, term.New("1"))))
	require.Equal(t, "2", mustFunc(Mul(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "4.08", mustFunc(Mul(nil, term.New("1.2"), term.New("3.4"))))
	require.Equal(t, "-22.848", mustFunc(Mul(nil, term.New("1.2"), term.New("3.4"), term.New("-5.6"))))
	require.Error(t, getError(Mul(nil, term.New("1.2"), term.New("3.4"), term.New("A"))))
}

func TestDiv(t *testing.T) {
	require.Error(t, getError(Div(nil)))
	require.Error(t, getError(Div(nil, term.New("1"))))
	require.Equal(t, "0.5", mustFunc(Div(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "0.3529411764705882", mustFunc(Div(nil, term.New("1.2"), term.New("3.4"))))
	require.Equal(t, "-0.0630252100840336", mustFunc(Div(nil, term.New("1.2"), term.New("3.4"), term.New("-5.6"))))
	require.Error(t, getError(Div(nil, term.New("1.2"), term.New("3.4"), term.New("A"))))
}

func TestMod(t *testing.T) {
	require.Error(t, getError(Mod(nil)))
	require.Error(t, getError(Mod(nil, term.New("1"))))
	require.Equal(t, "0", mustFunc(Mod(nil, term.New("2"), term.New("1"))))
	require.Equal(t, "1", mustFunc(Mod(nil, term.New("3"), term.New("2"))))
	require.Equal(t, "1", mustFunc(Mod(nil, term.New("4"), term.New("3"), term.New("2"))))
	require.Error(t, getError(Mod(nil, term.New("4"), term.New("3"), term.New("A"))))
}
