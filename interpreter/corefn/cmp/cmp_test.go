package cmp

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

func TestEqual(t *testing.T) {
	require.Error(t, getError(Equal(nil)))
	require.Error(t, getError(Equal(nil, term.New("1"))))
	require.Equal(t, "false", mustFunc(Equal(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "true", mustFunc(Equal(nil, term.New("1"), term.New("1"), term.New("1"))))
}

func TestNotEqual(t *testing.T) {
	require.Error(t, getError(NotEqual(nil)))
	require.Error(t, getError(NotEqual(nil, term.New("1"))))
	require.Equal(t, "true", mustFunc(NotEqual(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "false", mustFunc(NotEqual(nil, term.New("1"), term.New("1"), term.New("1"))))
}

func TestGreaterThan(t *testing.T) {
	require.Error(t, getError(GreaterThan(nil)))
	require.Error(t, getError(GreaterThan(nil, term.New("1"))))
	require.Equal(t, "true", mustFunc(GreaterThan(nil, term.New("2"), term.New("1"))))
	require.Equal(t, "true", mustFunc(GreaterThan(nil, term.New("3"), term.New("1"), term.New("2"))))
	require.Equal(t, "false", mustFunc(GreaterThan(nil, term.New("0"), term.New("-1"), term.New("2"))))
	require.Equal(t, "false", mustFunc(GreaterThan(nil, term.New("1"), term.New("1"))))
	require.Error(t, getError(GreaterThan(nil, term.New("1"), term.New("A"))))
}

func TestLessThan(t *testing.T) {
	require.Error(t, getError(LessThan(nil)))
	require.Error(t, getError(LessThan(nil, term.New("1"))))
	require.Equal(t, "true", mustFunc(LessThan(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "true", mustFunc(LessThan(nil, term.New("1"), term.New("2"), term.New("3"))))
	require.Equal(t, "false", mustFunc(LessThan(nil, term.New("1"), term.New("0"), term.New("2"))))
	require.Equal(t, "false", mustFunc(LessThan(nil, term.New("1"), term.New("1"))))
	require.Error(t, getError(LessThan(nil, term.New("1"), term.New("A"))))
}

func TestGreaterThanOrEqual(t *testing.T) {
	require.Error(t, getError(GreaterThanOrEqual(nil)))
	require.Error(t, getError(GreaterThanOrEqual(nil, term.New("1"))))
	require.Equal(t, "true", mustFunc(GreaterThanOrEqual(nil, term.New("2"), term.New("1"))))
	require.Equal(t, "true", mustFunc(GreaterThanOrEqual(nil, term.New("3"), term.New("1"), term.New("2"), term.New("3"))))
	require.Equal(t, "false", mustFunc(GreaterThanOrEqual(nil, term.New("0"), term.New("-1"), term.New("2"))))
	require.Error(t, getError(GreaterThanOrEqual(nil, term.New("1"), term.New("A"))))
}

func TestLessThanOrEqual(t *testing.T) {
	require.Error(t, getError(LessThanOrEqual(nil)))
	require.Error(t, getError(LessThanOrEqual(nil, term.New("1"))))
	require.Equal(t, "true", mustFunc(LessThanOrEqual(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "true", mustFunc(LessThanOrEqual(nil, term.New("1"), term.New("2"), term.New("3"), term.New("1"))))
	require.Equal(t, "false", mustFunc(LessThanOrEqual(nil, term.New("1"), term.New("0"), term.New("2"))))
	require.Error(t, getError(LessThanOrEqual(nil, term.New("1"), term.New("A"))))
}

func TestBetween(t *testing.T) {
	require.Error(t, getError(Between(nil)))
	require.Error(t, getError(Between(nil, term.New("1"))))
	require.Error(t, getError(Between(nil, term.New("1"), term.New("2"))))
	require.Equal(t, "true", mustFunc(Between(nil, term.New("1"), term.New("0"), term.New("2"))))
	require.Equal(t, "true", mustFunc(Between(nil, term.New("1"), term.New("1"), term.New("2"))))
	require.Equal(t, "true", mustFunc(Between(nil, term.New("1"), term.New("0"), term.New("1"))))
	require.Equal(t, "true", mustFunc(Between(nil, term.New("0"), term.New("1"), term.New("2"), term.New("0"), term.New("2"))))
	require.Equal(t, "false", mustFunc(Between(nil, term.New("3"), term.New("0"), term.New("2"))))
	require.Error(t, getError(Between(nil, term.New("1"), term.New("1"), term.New("A"))))
	require.Error(t, getError(Between(nil, term.New("1"), term.New("A"), term.New("1"))))
	require.Error(t, getError(Between(nil, term.New("A"), term.New("1"), term.New("1"))))
}
