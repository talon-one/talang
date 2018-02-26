package list_test

import (
	"log"
	"os"
	"testing"

	"github.com/talon-one/talang/lexer"

	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/interpreter/shared"

	"github.com/talon-one/talang/interpreter/corefn/list"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang/block"
)

func mustNewInterpreterWithLogger() *interpreter.Interpreter {
	interp := interpreter.MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)
	return interp
}

func mustFunc(result *block.Block, err error) *block.Block {
	if err != nil {
		panic(err)
	}
	return result
}

func getError(result interface{}, err error) error {
	return err
}

type Error struct{}

func (Error) Error() string { return "" }

type test struct {
	Input    string
	Binding  map[string]shared.Binding
	Expected interface{}
}

func runTests(t *testing.T, tests ...test) {
	interp := mustNewInterpreterWithLogger()
	require.NoError(t, interp.RemoveAllFunctions())
	require.NoError(t, interp.RegisterFunction(list.AllOperations()...))
	for i, test := range tests {
		interp := interp.NewScope()
		interp.Binding = test.Binding
		result, err := interp.LexAndEvaluate(test.Input)
		switch b := test.Expected.(type) {
		case error:
			require.Error(t, err)
		case *block.Block:
			require.EqualValues(t, test.Expected, result, "Test #%d failed, Expected %s was %s", i, b.Stringify(), result.Stringify())
		}
	}
}

func TestList(t *testing.T) {
	runTests(t, test{
		"list Hello World",
		nil,
		block.NewList(block.NewString("Hello"), block.NewString("World")),
	})
}

func TestHead(t *testing.T) {
	runTests(t,
		test{
			"head (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("Hello"),
		},
		test{
			"head (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewString("Hello"),
		},
		test{
			"head (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(),
				},
			},
			block.NewNull(),
		},
	)
}

func TestTail(t *testing.T) {
	runTests(t,
		test{
			"tail (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewList(block.NewString("World")),
		},
		test{
			"tail (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewList(),
		},
		test{
			"tail (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(),
				},
			},
			block.NewList(),
		},
	)
}

func TestDrop(t *testing.T) {
	runTests(t,
		test{
			"drop (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewList(block.NewString("Hello")),
		},
		test{
			"drop (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewList(),
		},
		test{
			"drop (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(),
				},
			},
			block.NewList(),
		},
	)
}

func TestItem(t *testing.T) {
	runTests(t,
		test{
			"item (. List) 0",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("Hello"),
		},
		test{
			"item (. List) 1",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("World"),
		},
		test{
			"item (. List) -1",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			Error{},
		},
		test{
			"item (. List) 2",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			Error{},
		},
		test{
			"item (. List) A",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			lexer.MustLex("item (. List) A"),
		},
	)
}
