package list_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/lexer"

	helpers "github.com/talon-one/talang/testhelpers"
	"github.com/talon-one/talang/token"
)

func TestList(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"list Hello World",
			nil,
			token.NewList(token.NewString("Hello"), token.NewString("World")),
		},
		helpers.Test{
			`list "Hello World" "Hello Universe"`,
			nil,
			token.NewList(token.NewString("Hello World"), token.NewString("Hello Universe")),
		},
	)
}

func TestHead(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"head (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
			}),
			token.NewString("Hello"),
		},
		helpers.Test{
			"head (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello")),
			}),
			token.NewString("Hello"),
		},
		helpers.Test{
			"head (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(),
			}),
			token.NewNull(),
		},
	)
}

func TestTail(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"tail (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
			}),
			token.NewList(token.NewString("World")),
		},
		helpers.Test{
			"tail (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello")),
			}),
			token.NewList(),
		},
		helpers.Test{
			"tail (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(),
			}),
			token.NewList(),
		},
	)
}

func TestDrop(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"drop (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
			}),
			token.NewList(token.NewString("Hello")),
		},
		helpers.Test{
			"drop (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello")),
			}),
			token.NewList(),
		},
		helpers.Test{
			"drop (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(),
			}),
			token.NewList(),
		},
	)
}

func TestItem(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"item (. List) 0",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
			}),
			token.NewString("Hello"),
		},
		helpers.Test{
			"item (. List) 1",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
			}),
			token.NewString("World"),
		},
		helpers.Test{
			"item (. List) -1",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
			}),
			helpers.Error{},
		},
		helpers.Test{
			"item (. List) 2",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
			}),
			helpers.Error{},
		},
		helpers.Test{
			"item (. List) A",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
			}),
			lexer.MustLex("item (. List) A"),
		},
	)
}

func TestPush(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.Binding = token.NewMap(map[string]*token.TaToken{
		"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
	})
	require.NoError(t, interp.RegisterTemplate(interpreter.TaTemplate{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn",
			Arguments: []token.Kind{
				token.String,
			},
			Returns: token.String,
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

	interp.Binding.MapItem("List").Children[0] = token.NewString("Dude!")
	require.EqualValues(t, interp.MustLexAndEvaluate("list Dude! World"), interp.Binding.MapItem("List"))
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World and"), newList)
}

func TestAppend(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.Binding = token.NewMap(map[string]*token.TaToken{
		"List": token.NewList(token.NewString("Hello"), token.NewString("World")),
	})
	require.NoError(t, interp.RegisterTemplate(interpreter.TaTemplate{
		CommonSignature: interpreter.CommonSignature{
			Name: "fn",
			Arguments: []token.Kind{
				token.String,
			},
			Returns: token.String,
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

	interp.Binding.MapItem("List").Children[0] = token.NewString("Dude!")
	require.EqualValues(t, interp.MustLexAndEvaluate("list Dude! World"), interp.Binding.MapItem("List"))
	require.EqualValues(t, interp.MustLexAndEvaluate("list Hello World and"), newList)
}

func TestMap(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`map (. List) x (+ (. x Name) " " (. x Surname))`,
		token.NewMap(map[string]*token.TaToken{
			"List": token.NewList(
				token.NewMap(map[string]*token.TaToken{
					"Name":    token.NewString("Joe"),
					"Surname": token.NewString("Doe"),
					"Id":      token.NewDecimalFromInt(0),
				}),
				token.NewMap(map[string]*token.TaToken{
					"Name":    token.NewString("Alice"),
					"Surname": token.NewString("Wonder"),
					"Id":      token.NewDecimalFromInt(1),
				}),
			),
		}),
		token.NewList(token.NewString("Joe Doe"), token.NewString("Alice Wonder")),
	})
}

func TestMapLegacy(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		`map (. List) ((x) (+ (. x Name) " " (. x Surname)))`,
		token.NewMap(map[string]*token.TaToken{
			"List": token.NewList(
				token.NewMap(map[string]*token.TaToken{
					"Name":    token.NewString("Joe"),
					"Surname": token.NewString("Doe"),
					"Id":      token.NewDecimalFromInt(0),
				}),
				token.NewMap(map[string]*token.TaToken{
					"Name":    token.NewString("Alice"),
					"Surname": token.NewString("Wonder"),
					"Id":      token.NewDecimalFromInt(1),
				}),
			),
		}),
		token.NewList(token.NewString("Joe Doe"), token.NewString("Alice Wonder")),
	}, helpers.Test{
		"map (list 1 2 3) (4)",
		nil,
		helpers.Error{},
	})
}

func TestSort(t *testing.T) {
	interp := helpers.MustNewInterpreterWithLogger()
	interp.Binding = token.NewMap(map[string]*token.TaToken{
		"List": token.NewList(
			token.NewString("World"),
			token.NewDecimalFromInt(2),
			token.NewString("Hello"),
			token.NewDecimalFromInt(1),
		),
	})
	require.EqualValues(t, token.NewList(token.NewDecimalFromInt(1), token.NewDecimalFromInt(2), token.NewString("Hello"), token.NewString("World")), interp.MustLexAndEvaluate("sort (. List)"))

	// integrity check
	require.EqualValues(t, token.NewList(token.NewString("World"), token.NewDecimalFromInt(2), token.NewString("Hello"), token.NewDecimalFromInt(1)), interp.Get("List"))

	require.EqualValues(t, token.NewList(token.NewString("World"), token.NewString("Hello"), token.NewDecimalFromInt(2), token.NewDecimalFromInt(1)), interp.MustLexAndEvaluate("sort (. List) true"))

	// integrity check
	require.EqualValues(t, token.NewList(token.NewString("World"), token.NewDecimalFromInt(2), token.NewString("Hello"), token.NewDecimalFromInt(1)), interp.Get("List"))
}

func TestMin(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`min (list 100 4 3 10 6000 90 99)`,
			nil,
			token.NewDecimalFromInt(3),
		},
		helpers.Test{
			`min (list 100 4 3 10 6000 Hello 90 99)`,
			nil,
			token.NewDecimalFromInt(3),
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
			token.NewDecimalFromInt(6000),
		},
		helpers.Test{
			`max (list 100 4 3 10 6000 Hello 90 99)`,
			nil,
			token.NewDecimalFromInt(6000),
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
			token.NewDecimalFromInt(4),
		},
		helpers.Test{
			`count (list hola hola amigos)`,
			nil,
			token.NewDecimalFromInt(3),
		},
	)
}

func TestReverse(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`reverse (list 1 2 3 4)`,
			nil,
			token.NewList(
				token.NewDecimalFromInt(4),
				token.NewDecimalFromInt(3),
				token.NewDecimalFromInt(2),
				token.NewDecimalFromInt(1),
			),
		},
	)
}

func BenchmarkReverse(b *testing.B) {
	tests := []struct {
		input    string
		expected *token.TaToken
	}{
		{"reverse (list 1 2 3 4)", token.NewList(
			token.NewDecimalFromInt(4),
			token.NewDecimalFromInt(3),
			token.NewDecimalFromInt(2),
			token.NewDecimalFromInt(1),
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
			token.NewString("hello-world"),
		},
		helpers.Test{
			`join (. lol) -`,
			token.NewMap(map[string]*token.TaToken{
				"lol": token.NewList(token.NewString("lol")),
			}),
			token.NewString("lol"),
		},
		helpers.Test{
			`join (. lol) -`,
			token.NewMap(map[string]*token.TaToken{
				"lol": token.NewList(token.NewString("lo ll")),
			}),
			token.NewString("lo ll"),
		},
		helpers.Test{
			`join (list lo        l) -`,
			nil,
			token.NewString("lo-l"),
		},
		helpers.Test{
			`join (list 1 2 3) -`,
			nil,
			helpers.Error{},
		},
		helpers.Test{
			`join (list "1" "2" "3") -`,
			nil,
			token.NewString("1-2-3"),
		},
	)
}

func TestSplit(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`split "1,2,3,a" ","`,
			nil,
			token.NewList(
				token.NewString("1"),
				token.NewString("2"),
				token.NewString("3"),
				token.NewString("a"),
			),
		},
	)
}
func TestIsEmpty(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`isEmpty (list hello world)`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`isEmpty (list)`,
			nil,
			token.NewBool(true),
		},
	)
}

func TestExists(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`exists (list hello world) Item (= (. Item) "hello")`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`exists (list hello world) Item (= (. Item) "world")`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`exists (list hello world) Item (= (. Item) "universes")`,
			nil,
			token.NewBool(false),
		}, helpers.Test{
			`exists (list hello world) Item (+ (. Item) "world")`,
			nil,
			helpers.Error{},
		}, helpers.Test{
			`exists (list hello world) Item (panic)`,
			nil,
			helpers.Error{},
		},
	)
}

func TestExistsLegacy(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`exists (list hello world) ((Item) (= (. Item) "hello"))`,
			nil,
			token.NewBool(true),
		}, helpers.Test{
			`exists (list hello world) (panic)`,
			nil,
			helpers.Error{},
		},
	)
}

func TestSum(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`sum (. List) Item (. Item Price)`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(2),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(2),
					}),
				),
			}),
			token.NewDecimalFromInt(4),
		}, helpers.Test{
			`sum (. List) Item (panic)`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(2),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(2),
					}),
				),
			}),
			helpers.Error{},
		}, helpers.Test{
			`sum (. List) Item (. Item Price)`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewString("Hey"),
					}),
				),
			}),
			helpers.Error{},
		},
	)
}

func TestEvery(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`every (. List) Item (= (. Item Price) 1)`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(1),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(1),
					}),
				),
			}),
			token.NewBool(true),
		}, helpers.Test{
			`every (. List) Item (= (. Item Price) 1)`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(2),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(1),
					}),
				),
			}),
			token.NewBool(false),
		}, helpers.Test{
			`every (. List) Item (= (panic) 1)`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(2),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(1),
					}),
				),
			}),
			helpers.Error{},
		}, helpers.Test{
			`every (. List) Item (= (. Item SKU) "XJK_992")`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"SKU": token.NewString("XJK_992"),
					}),
					token.NewMap(map[string]*token.TaToken{
						"SKU": token.NewString("XJK_992"),
					}),
				),
			}),
			token.NewBool(true),
		},
	)
}

func TestEveryLegacy(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`every (. List) ((Item) (= (. Item Price) 1))`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(1),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(1),
					}),
				),
			}),
			token.NewBool(true),
		}, helpers.Test{
			`every (. List) ((Item) (= (. Item Price) 1))`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(2),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(1),
					}),
				),
			}),
			token.NewBool(false),
		}, helpers.Test{
			`every (. List) (panic)`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(2),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Price": token.NewDecimalFromInt(1),
					}),
				),
			}),
			helpers.Error{},
		}, helpers.Test{
			`every (. List) ((Item) (= (. Item SKU) "XJK_992"))`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"SKU": token.NewString("XJK_992"),
					}),
					token.NewMap(map[string]*token.TaToken{
						"SKU": token.NewString("XJK_992"),
					}),
				),
			}),
			token.NewBool(true),
		},
	)
}

func TestSortByNumber(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`sortByNumber (. List) ((Item) (. Item Price)) false`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Name":  token.NewString("Alex"),
						"Price": token.NewDecimalFromInt(26),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Name":  token.NewString("Gertrude"),
						"Price": token.NewDecimalFromInt(11),
					}),
				),
			}),
			token.NewList(
				token.NewMap(map[string]*token.TaToken{
					"Name":  token.NewString("Gertrude"),
					"Price": token.NewDecimalFromInt(11),
				}),
				token.NewMap(map[string]*token.TaToken{
					"Name":  token.NewString("Alex"),
					"Price": token.NewDecimalFromInt(26),
				}),
			),
		}, helpers.Test{
			`sortByNumber (. List) ((Item) (. Item Price)) true`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Name":  token.NewString("Alex"),
						"Price": token.NewDecimalFromInt(26),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Name":  token.NewString("Gertrude"),
						"Price": token.NewDecimalFromInt(11),
					}),
				),
			}),
			token.NewList(
				token.NewMap(map[string]*token.TaToken{
					"Name":  token.NewString("Alex"),
					"Price": token.NewDecimalFromInt(26),
				}),
				token.NewMap(map[string]*token.TaToken{
					"Name":  token.NewString("Gertrude"),
					"Price": token.NewDecimalFromInt(11),
				}),
			),
		}, helpers.Test{
			`sortByNumber (list 2 4 3 1) ((Item) (. Item)) true`,
			nil,
			token.NewList(
				token.NewDecimalFromInt(4),
				token.NewDecimalFromInt(3),
				token.NewDecimalFromInt(2),
				token.NewDecimalFromInt(1),
			),
		}, helpers.Test{
			`sortByNumber (list 2 4 3 1) ((Item) (. Item)) false`,
			nil,
			token.NewList(
				token.NewDecimalFromInt(1),
				token.NewDecimalFromInt(2),
				token.NewDecimalFromInt(3),
				token.NewDecimalFromInt(4),
			),
		},
	)
}

func TestSortByString(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`sortByString (. List) ((Item) (. Item Name)) false`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Name":  token.NewString("Gertrude"),
						"Price": token.NewDecimalFromInt(26),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Name":  token.NewString("Alex"),
						"Price": token.NewDecimalFromInt(11),
					}),
				),
			}),
			token.NewList(
				token.NewMap(map[string]*token.TaToken{
					"Name":  token.NewString("Alex"),
					"Price": token.NewDecimalFromInt(11),
				}),
				token.NewMap(map[string]*token.TaToken{
					"Name":  token.NewString("Gertrude"),
					"Price": token.NewDecimalFromInt(26),
				}),
			),
		}, helpers.Test{
			`sortByString (. List) ((Item) (. Item Name)) true`,
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewList(
					token.NewMap(map[string]*token.TaToken{
						"Name":  token.NewString("Alex"),
						"Price": token.NewDecimalFromInt(26),
					}),
					token.NewMap(map[string]*token.TaToken{
						"Name":  token.NewString("Gertrude"),
						"Price": token.NewDecimalFromInt(11),
					}),
				),
			}),
			token.NewList(
				token.NewMap(map[string]*token.TaToken{
					"Name":  token.NewString("Gertrude"),
					"Price": token.NewDecimalFromInt(11),
				}),
				token.NewMap(map[string]*token.TaToken{
					"Name":  token.NewString("Alex"),
					"Price": token.NewDecimalFromInt(26),
				}),
			),
		}, helpers.Test{
			`sortByString (list "b" "a" "z" "t") ((Item) (. Item)) false`,
			nil,
			token.NewList(
				token.NewString("a"),
				token.NewString("b"),
				token.NewString("t"),
				token.NewString("z"),
			),
		}, helpers.Test{
			`sortByString (list "b" "a" "z" "t") ((Item) (. Item)) true`,
			nil,
			token.NewList(
				token.NewString("z"),
				token.NewString("t"),
				token.NewString("b"),
				token.NewString("a"),
			),
		},
	)
}
