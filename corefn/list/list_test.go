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

	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World"), interp.Binding.MapItem("List"))

	newList := interp.MustLexAndEvaluate("push (. List) and")

	interp.Binding.MapItem("List").Children[0] = block.NewString("Dude!")
	require.EqualValues(t, interp.MustLexAndEvaluate("list Dude! World"), interp.Binding.MapItem("List"))
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World and"), newList)
}

func TestAppend(t *testing.T) {
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
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World and Universe"), interp.MustLexAndEvaluate("append (. List) and Universe"))

	// check if the original list is still unmodified
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World"), interp.Binding.MapItem("List"))

	// Push with a function inside
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World Alice"), interp.MustLexAndEvaluate("append (. List) (! fn Alice)"))

	// check if the original list is still unmodified
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World"), interp.Binding.MapItem("List"))

	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World"), interp.Binding.MapItem("List"))

	newList := interp.MustLexAndEvaluate("append (. List) and")

	interp.Binding.MapItem("List").Children[0] = block.NewString("Dude!")
	require.EqualValues(t, interp.MustLexAndEvaluate("list Dude! World"), interp.Binding.MapItem("List"))
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World and"), newList)
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

func TestSort(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.Binding = block.NewMap(map[string]*block.Block{
		"List": block.NewList(
			block.NewString("World"),
			block.NewDecimalFromInt(2),
			block.NewString("Hello"),
			block.NewDecimalFromInt(1),
		),
	})
	require.EqualValues(t, block.NewList(block.NewDecimalFromInt(1), block.NewDecimalFromInt(2), block.NewString("Hello"), block.NewString("World")), interp.MustLexAndEvaluate("sort (. List)"))

	// integrity check
	require.EqualValues(t, block.NewList(block.NewString("World"), block.NewDecimalFromInt(2), block.NewString("Hello"), block.NewDecimalFromInt(1)), interp.Get("List"))

	require.EqualValues(t, block.NewList(block.NewString("World"), block.NewString("Hello"), block.NewDecimalFromInt(2), block.NewDecimalFromInt(1)), interp.MustLexAndEvaluate("sort (. List) true"))

	// integrity check
	require.EqualValues(t, block.NewList(block.NewString("World"), block.NewDecimalFromInt(2), block.NewString("Hello"), block.NewDecimalFromInt(1)), interp.Get("List"))
}

func TestMin(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`min (list 100 4 3 10 6000 90 99)`,
			nil,
			block.NewDecimalFromInt(3),
		},
		helpers.Test{
			`min (list 100 4 3 10 6000 Hello 90 99)`,
			nil,
			block.NewDecimalFromInt(3),
		},
		helpers.Test{
			`min (list Hello World)`,
			nil,
			&helpers.Error{},
		},
	)
}

func TestMax(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`max (list 100 4 3 10 6000 90 99)`,
			nil,
			block.NewDecimalFromInt(6000),
		},
		helpers.Test{
			`max (list 100 4 3 10 6000 Hello 90 99)`,
			nil,
			block.NewDecimalFromInt(6000),
		},
		helpers.Test{
			`max (list Hello World)`,
			nil,
			&helpers.Error{},
		},
	)
}

func TestCount(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`count (list 1 2 3 4)`,
			nil,
			block.NewDecimalFromInt(4),
		},
		helpers.Test{
			`count (list hola hola amigos)`,
			nil,
			block.NewDecimalFromInt(3),
		},
	)
}

func TestReverse(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`reverse (list 1 2 3 4)`,
			nil,
			block.NewList(
				block.NewDecimalFromInt(4),
				block.NewDecimalFromInt(3),
				block.NewDecimalFromInt(2),
				block.NewDecimalFromInt(1),
			),
		},
	)
}

func BenchmarkReverse(b *testing.B) {
	tests := []struct {
		input    string
		expected *block.Block
	}{
		{"reverse (list 1 2 3 4)", block.NewList(
			block.NewDecimalFromInt(4),
			block.NewDecimalFromInt(3),
			block.NewDecimalFromInt(2),
			block.NewDecimalFromInt(1),
		)},
	}

	interp := helpers.MustNewInterpreter()

	for i := 0; i < b.N; i++ {
		for _, test := range tests {
			require.Equal(b, test.expected, interp.MustLexAndEvaluate(test.input))
		}
	}
}

func TestJoin(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`join (list hello world) -`,
			nil,
			block.NewString("hello-world"),
		},
		helpers.Test{
			`join (. lol) -`,
			block.NewMap(map[string]*block.Block{
				"lol": block.NewList(block.NewString("lol")),
			}),
			block.NewString("lol"),
		},
		helpers.Test{
			`join (. lol) -`,
			block.NewMap(map[string]*block.Block{
				"lol": block.NewList(block.NewString("lo ll")),
			}),
			block.NewString("lo ll"),
		},
		helpers.Test{
			`join (list lo        l) -`,
			nil,
			block.NewString("lo-l"),
		},
		helpers.Test{
			`join (list 1 2 3) -`,
			nil,
			helpers.Error{},
		},
		// FIX - the interpreter parses "1" as a decimal, should be a string.
		// helpers.Test{
		// 	`join (list "1" "2" "3") -`,
		// 	nil,
		// 	block.NewString("1-2-3"),
		// },
	)
}

func TestIsEmpty(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`isEmpty (list hello world)`,
			nil,
			block.NewBool(false),
		},
	)
}
