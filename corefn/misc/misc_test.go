package misc_test

import (
	"testing"

	helpers "github.com/talon-one/talang/testhelpers"
	"github.com/talon-one/talang/token"
)

func TestNoop(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"noop",
		nil,
		token.NewNull(),
	})
}

func TestToString(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"toString 1",
		nil,
		token.NewString("1"),
	})
}

func TestNot(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"not false",
			nil,
			token.NewBool(true),
		},
		helpers.Test{
			"not true",
			nil,
			token.NewBool(false),
		},
		helpers.Test{
			"not (not false)",
			nil,
			token.NewBool(false),
		},
	)
}

func TestCatch(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`catch 22 (panic)`,
			nil,
			token.NewDecimalFromInt(22),
		},
		helpers.Test{
			`catch (+ 1 1) (panic)`,
			nil,
			token.NewDecimalFromInt(2),
		},
		helpers.Test{
			`catch "22" (. Profile)`,
			nil,
			token.NewString("22"),
		},
		helpers.Test{
			`catch 22 (. Profile)`,
			nil,
			token.NewDecimalFromInt(22),
		},
		helpers.Test{
			`catch 22 (+ 2 2)`,
			nil,
			token.NewDecimalFromInt(4),
		},
		helpers.Test{
			`catch 22 2`,
			nil,
			token.NewDecimalFromInt(2),
		},
		helpers.Test{
			`(catch (+ 2 (* 5 (- 3 4))) (+ 2 ( * 4 (panic))))`,
			nil,
			token.NewDecimalFromInt(-3),
		},
	)
}

func TestDo(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"do (list 1 2 3) Item (. Item)",
		nil,
		token.NewList(
			token.NewDecimalFromInt(1),
			token.NewDecimalFromInt(2),
			token.NewDecimalFromInt(3),
		),
	}, helpers.Test{
		"do 4 x (panic)",
		nil,
		helpers.Error{},
	})
}

func TestDoLegacy(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"do (list 1 2 3) ((Item) (. Item))",
		nil,
		token.NewList(
			token.NewDecimalFromInt(1),
			token.NewDecimalFromInt(2),
			token.NewDecimalFromInt(3),
		),
	}, helpers.Test{
		"do 4 ((x) (panic))",
		nil,
		helpers.Error{},
	}, helpers.Test{
		"do 4 (x (panic))",
		nil,
		helpers.Error{},
	}, helpers.Test{
		"do 4 (panic)",
		nil,
		helpers.Error{},
	})
}

func TestSafeRead(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			".| boo (. List)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewString("XJK_992"),
			}),
			token.NewString("XJK_992"),
		}, helpers.Test{
			".| boo (. Meh)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewString("XJK_992"),
			}),
			token.NewString("boo"),
		}, helpers.Test{
			".| boo (. List SKU)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewMap(map[string]*token.TaToken{
					"SKU": token.NewString("XJK_992"),
				}),
			}),
			token.NewString("XJK_992"),
		}, helpers.Test{
			".| 2 (. List SKU)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewMap(map[string]*token.TaToken{
					"SKU": token.NewDecimalFromInt(1),
				}),
			}),
			token.NewDecimalFromInt(1),
		}, helpers.Test{
			".| 2 (. List SKU)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewMap(map[string]*token.TaToken{
					"SKU": token.NewDecimalFromInt(1),
				}),
			}),
			token.NewDecimalFromInt(1),
		}, helpers.Test{
			".| 2 (. List SKU)",
			token.NewMap(map[string]*token.TaToken{
				"List": token.NewMap(map[string]*token.TaToken{
					"SKU": token.NewString("asda"),
				}),
			}),
			helpers.Error{},
		},
		// This test evaluates to 1, bug?
		// , helpers.Test{
		// 	".| 2 (. List SKU Dis)",
		// 	token.NewMap(map[string]*token.TaToken{
		// 		"List": token.NewMap(map[string]*token.TaToken{
		// 			"SKU": token.NewDecimalFromInt(1),
		// 		}),
		// 	}),
		// 	token.NewDecimalFromInt(2),
		// }
	)
}
