package list_test

import (
	"testing"

	"github.com/talon-one/talang/lexer"

	"github.com/talon-one/talang/interpreter/shared"

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
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("Hello"),
		},
		helpers.Test{
			"head (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewString("Hello"),
		},
		helpers.Test{
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
	helpers.RunTests(t,
		helpers.Test{
			"tail (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewList(block.NewString("World")),
		},
		helpers.Test{
			"tail (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewList(),
		},
		helpers.Test{
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
	helpers.RunTests(t,
		helpers.Test{
			"drop (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewList(block.NewString("Hello")),
		},
		helpers.Test{
			"drop (. List)",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello")),
				},
			},
			block.NewList(),
		},
		helpers.Test{
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
	helpers.RunTests(t,
		helpers.Test{
			"item (. List) 0",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("Hello"),
		},
		helpers.Test{
			"item (. List) 1",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			block.NewString("World"),
		},
		helpers.Test{
			"item (. List) -1",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			helpers.Error{},
		},
		helpers.Test{
			"item (. List) 2",
			map[string]shared.Binding{
				"List": shared.Binding{
					Value: block.NewList(block.NewString("Hello"), block.NewString("World")),
				},
			},
			helpers.Error{},
		},
		helpers.Test{
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
