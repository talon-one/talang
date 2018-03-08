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
	helpers.RunTests(t,
		helpers.Test{
			"list Hello World",
			nil,
			block.NewList(block.NewString("Hello"), block.NewString("World")),
		},
		helpers.Test{
			`list "Hello World" "Hello Universe"`,
			nil,
			block.NewList(block.NewString("Hello World"), block.NewString("Hello Universe")),
		},
	)
}

func TestHead(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"head (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
			}),
			block.NewString("Hello"),
		},
		helpers.Test{
			"head (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello")),
			}),
			block.NewString("Hello"),
		},
		helpers.Test{
			"head (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(),
			}),
			block.NewNull(),
		},
	)
}

func TestTail(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"tail (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
			}),
			block.NewList(block.NewString("World")),
		},
		helpers.Test{
			"tail (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello")),
			}),
			block.NewList(),
		},
		helpers.Test{
			"tail (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(),
			}),
			block.NewList(),
		},
	)
}

func TestDrop(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"drop (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
			}),
			block.NewList(block.NewString("Hello")),
		},
		helpers.Test{
			"drop (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello")),
			}),
			block.NewList(),
		},
		helpers.Test{
			"drop (. List)",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(),
			}),
			block.NewList(),
		},
	)
}

func TestItem(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"item (. List) 0",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
			}),
			block.NewString("Hello"),
		},
		helpers.Test{
			"item (. List) 1",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
			}),
			block.NewString("World"),
		},
		helpers.Test{
			"item (. List) -1",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
			}),
			helpers.Error{},
		},
		helpers.Test{
			"item (. List) 2",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
			}),
			helpers.Error{},
		},
		helpers.Test{
			"item (. List) A",
			block.NewMap(map[string]*block.Block{
				"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
			}),
			lexer.MustLex("item (. List) A"),
		},
	)
}

func TestPush(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.Binding = block.NewMap(map[string]*block.Block{
		"List": block.NewList(block.NewString("Hello"), block.NewString("World")),
	})
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
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World"), interp.Binding.MapItem("List"))

	// Push with a function inside
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World Alice"), interp.MustLexAndEvaluate("push (. List) (! fn Alice)"))

	// check if the original list is still unmodified
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World"), interp.Binding.MapItem("List"))
}

func TestMap(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`map (. List) x (+ (. x Name) " " (. x Surname))`,
		block.NewMap(map[string]*block.Block{
			"List": block.NewList(
				block.NewMap(map[string]*block.Block{
					"Name":    block.NewString("Joe"),
					"Surname": block.NewString("Doe"),
					"Id":      block.NewDecimalFromInt(0),
				}),
				block.NewMap(map[string]*block.Block{
					"Name":    block.NewString("Alice"),
					"Surname": block.NewString("Wonder"),
					"Id":      block.NewDecimalFromInt(1),
				}),
			),
		}),
		block.NewList(block.NewString("Joe Doe"), block.NewString("Alice Wonder")),
	})
}
