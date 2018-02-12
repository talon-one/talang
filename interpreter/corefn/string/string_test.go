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

func TestStartsWith(t *testing.T) {
	require.Error(t, getError(StartsWith.Func(nil, []*block.Block{})))
	require.Error(t, getError(StartsWith.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "false", mustFunc(StartsWith.Func(nil, []*block.Block{block.New("Hello World"), block.New("Bye")})))
	require.Equal(t, "false", mustFunc(StartsWith.Func(nil, []*block.Block{block.New("Hello World"), block.New("Hello"), block.New("Bye")})))
	require.Equal(t, "true", mustFunc(StartsWith.Func(nil, []*block.Block{block.New("Hello World"), block.New("Hello"), block.New("Hell")})))
}

func TestEndsWith(t *testing.T) {
	require.Error(t, getError(EndsWith.Func(nil, []*block.Block{})))
	require.Error(t, getError(EndsWith.Func(nil, []*block.Block{block.New("1")})))
	require.Equal(t, "false", mustFunc(EndsWith.Func(nil, []*block.Block{block.New("Hello World"), block.New("Universe")})))
	require.Equal(t, "false", mustFunc(EndsWith.Func(nil, []*block.Block{block.New("Hello World"), block.New("World"), block.New("Universe")})))
	require.Equal(t, "true", mustFunc(EndsWith.Func(nil, []*block.Block{block.New("Hello World"), block.New("World"), block.New("ld")})))
}

func TestRegexp(t *testing.T) {
	require.Error(t, getError(Regexp.Func(nil, []*block.Block{})))
	require.Error(t, getError(Regexp.Func(nil, []*block.Block{block.New("foo")})))
	require.Equal(t, "true", mustFunc(Regexp.Func(nil, []*block.Block{block.New("foobar"), block.New("^foo")})))
	require.Equal(t, "false", mustFunc(Regexp.Func(nil, []*block.Block{block.New("foobar"), block.New("^foo$")})))
}

func TestAllOperations(t *testing.T) {
	AllOperations()
}
