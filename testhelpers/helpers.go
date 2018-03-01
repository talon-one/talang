package testhelpers

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
)

func MustNewInterpreterWithLogger() *talang.Interpreter {
	interp := talang.MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)
	return interp
}

func MustBlock(result *block.Block, err error) *block.Block {
	if err != nil {
		panic(err)
	}
	return result
}

func MustError(result interface{}, err error) error {
	return err
}

type Error struct{}

func (Error) Error() string { return "" }

type Test struct {
	Input    string
	Binding  map[string]interpreter.Binding
	Expected interface{}
}

func RunTests(t *testing.T, tests ...Test) {
	RunTestsWithInterpreter(t, MustNewInterpreterWithLogger(), tests...)
}

func RunTestsWithInterpreter(t *testing.T, interp *talang.Interpreter, tests ...Test) {
	for i, test := range tests {
		interp := interp.NewScope()
		interp.Binding = test.Binding
		result, err := interp.LexAndEvaluate(test.Input)
		switch b := test.Expected.(type) {
		case error:
			require.Error(t, err, "Test #%d failed", i)
		case *block.Block:
			if b.IsDecimal() && result.IsDecimal() {
				require.Equal(t, 0, b.Decimal.Cmp(result.Decimal))
			} else {
				require.EqualValues(t, test.Expected, result, "Test #%d failed, Expected %s was %s", i, b.Stringify(), result.Stringify())
			}
		}
	}
}
