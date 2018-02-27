package mapping_test

import (
	"testing"

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
	)
}
