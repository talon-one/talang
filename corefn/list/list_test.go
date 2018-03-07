package list_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/lexer"

	"github.com/talon-one/talang/block"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestList(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"list Hello World",
		nil,
		block.NewList(block.NewString("Hello"), block.NewString("World")),
	})
}

func TestHead(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"head (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("Hello"),
		},
		helpers.Test{
			"head (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewString("Hello"),
		},
		helpers.Test{
			"head (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(),
				},
			},
			block.NewNull(),
		},
	)
}

func TestTail(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"tail (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewList(block.NewString("World")),
		},
		helpers.Test{
			"tail (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewList(),
		},
		helpers.Test{
			"tail (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(),
				},
			},
			block.NewList(),
		},
	)
}

func TestDrop(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"drop (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewList(block.NewString("Hello")),
		},
		helpers.Test{
			"drop (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewList(),
		},
		helpers.Test{
			"drop (. List)",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(),
				},
			},
			block.NewList(),
		},
	)
}

func TestItem(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"item (. List) 0",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("Hello"),
		},
		helpers.Test{
			"item (. List) 1",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("World"),
		},
		helpers.Test{
			"item (. List) -1",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			helpers.Error{},
		},
		helpers.Test{
			"item (. List) 2",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			helpers.Error{},
		},
		helpers.Test{
			"item (. List) A",
			map[string]interpreter.Binding{
				"List": interpreter.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			lexer.MustLex("item (. List) A"),
		},
	)
}

func TestPush(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.Binding = map[string]interpreter.Binding{
		"List": interpreter.Binding{
			Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
		},
	}
	require.NoError(t, interp.RegisterTemplate(interpreter.TaTemplate{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn",
			Arguments: []block.Kind{
				block.StringKind,
			},
			Returns: block.StringKind,
		},
		Template: *lexer.MustLex(`(# 0)`),
	}))
	// check if the return value contains the appended data
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World and Universe"), interp.MustLexAndEvaluate("push (. List) and Universe"))

	// check if the original list is still unmodified
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World"), interp.Binding["List"].Value)

	// Push with a function inside
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World Alice"), interp.MustLexAndEvaluate("push (. List) (! fn Alice)"))

	// check if the original list is still unmodified
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World"), interp.Binding["List"].Value)
}
