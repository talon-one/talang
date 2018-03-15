package misc_test

import (
	"errors"
	"testing"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
	helpers "github.com/talon-one/talang/testhelpers"
)

func TestNoop(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"noop",
		nil,
		block.NewNull(),
	})
}

func TestToString(t *testing.T) {
	helpers.RunTests(t, helpers.Test{
		"toString 1",
		nil,
		block.NewString("1"),
	})
}

func TestNot(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			"not false",
			nil,
			block.NewBool(true),
		},
		helpers.Test{
			"not true",
			nil,
			block.NewBool(false),
		},
		helpers.Test{
			"not (not false)",
			nil,
			block.NewBool(false),
		},
	)
}

func TestCatch(t *testing.T) {
	helpers.RunTests(t,
		helpers.Test{
			`catch 22 (panic)`,
			nil,
			block.NewDecimalFromInt(22),
		},
		helpers.Test{
			`catch (+ 1 1) (panic)`,
			nil,
			block.NewDecimalFromInt(2),
		},
		helpers.Test{
			`catch "22" (. Profile)`,
			nil,
			block.NewString("22"),
		},
		helpers.Test{
			`catch 22 (. Profile)`,
			nil,
			block.NewDecimalFromInt(22),
		},
		helpers.Test{
			`catch 22 (+ 2 2)`,
			nil,
			block.NewDecimalFromInt(4),
		},
		helpers.Test{
			`catch 22 2`,
			nil,
			block.NewDecimalFromInt(2),
		},
		helpers.Test{
			`(catch (+ 2 (* 5 (- 3 4))) (+ 2 ( * 4 (panic))))`,
			nil,
			block.NewDecimalFromInt(-3),
		},
	)
}

func init() {
	interpreter.RegisterCoreFunction(
		interpreter.TaFunction{
			CommonSignature: interpreter.CommonSignature{
				Name:    "panic",
				Returns: block.AnyKind,
			},
			Func: func(interp *interpreter.Interpreter, args ...*block.Block) (*block.Block, error) {
				return nil, errors.New("panic")
			},
		},
	)
}
