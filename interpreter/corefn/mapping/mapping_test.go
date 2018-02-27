package mapping_test

import (
	"testing"

	"github.com/ericlagergren/decimal"
	"github.com/talon-one/talang/block"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestKV(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"kv (Key1 Value1) (Key2 Value2 Value3)",
			nil,
			block.NewMap(map[string]*block.Block{
				"Key1": block.NewString("Value1"),
				"Key2": block.NewList(block.NewString("Value2"), block.NewString("Value3")),
			}),
		},
		helpers.Test{
			"kv (Key1 (+ 1 2))",
			nil,
			block.NewMap(map[string]*block.Block{
				"Key1": block.NewDecimal(decimal.New(3, 0)),
			}),
		},
		helpers.Test{
			"kv (Key1 (kv (SubKey Value)))",
			nil,
			block.NewMap(map[string]*block.Block{
				"Key1": block.NewMap(map[string]*block.Block{
					"SubKey": block.NewString("Value"),
				}),
			}),
		},
	)
}
