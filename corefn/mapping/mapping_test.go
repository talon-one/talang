package mapping_test

import (
	"testing"

	helpers "github.com/talon-one/talang/testhelpers"
	"github.com/talon-one/talang/token"
)

func TestKV(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"kv (Key1 Value1) (Key2 Value2)",
			nil,
			token.NewMap(map[string]*token.TaToken{
				"Key1": token.NewString("Value1"),
				"Key2": token.NewString("Value2"),
			}),
		},
		helpers.Test{
			"kv (Key1 Value1) (Key2 (list Value2 Value3))",
			nil,
			token.NewMap(map[string]*token.TaToken{
				"Key1": token.NewString("Value1"),
				"Key2": token.NewList(token.NewString("Value2"), token.NewString("Value3")),
			}),
		},
		helpers.Test{
			"kv (Key1 Value1) (Key2 Value2 Value3)",
			nil,
			&helpers.Error{},
		},
		helpers.Test{
			"kv (Key1 (+ 1 2))",
			nil,
			token.NewMap(map[string]*token.TaToken{
				"Key1": token.NewDecimalFromInt(3),
			}),
		},
		helpers.Test{
			"kv (Key1 (kv (SubKey Value)))",
			nil,
			token.NewMap(map[string]*token.TaToken{
				"Key1": token.NewMap(map[string]*token.TaToken{
					"SubKey": token.NewString("Value"),
				}),
			}),
		},
		helpers.Test{
			"kv (Key1 (panic))",
			nil,
			&helpers.Error{},
		},
	)
}
