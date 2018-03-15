package testhelpers

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/talon-one/talang"
	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/token"
)

func init() {
	interpreter.RegisterCoreFunction(
		interpreter.TaFunction{
			CommonSignature: interpreter.CommonSignature{
				Name:    "panic",
				Returns: token.Any,
			},
			Func: func(interp *interpreter.Interpreter, args ...*token.TaToken) (*token.TaToken, error) {
				return nil, errors.New("panic")
			},
		},
	)
}

func MustNewInterpreterWithLogger() *talang.Interpreter {
	interp := talang.MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", 0)
	return interp
}

func MustNewInterpreter() *talang.Interpreter {
	return talang.MustNewInterpreter()
}

func MustBlock(result *token.TaToken, err error) *token.TaToken {
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
	Binding  *token.TaToken
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
			require.Error(t, err, "Test %d failed", i)
		case *token.TaToken:
			require.EqualValues(t, true, b.Equal(result), "Test #%d failed, Expected %s was %s", i, b.Stringify(), result.Stringify())
		}
	}
}
