package interpreter

import (
	"strings"

	"log"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/interpreter/internal"
	lexer "github.com/talon-one/talang/lexer"
	"github.com/talon-one/talang/term"
)

type Interpreter struct {
	internal.Interpreter
	functionMap map[string]TaFunc
	Logger      *log.Logger
}

func NewInterpreter() (*Interpreter, error) {
	interp := Interpreter{
		functionMap: make(map[string]TaFunc),
	}
	if err := interp.registerCoreFunctions(); err != nil {
		return nil, err
	}
	interp.Logger = log.New(&dummyWriter{}, "", 0)
	return &interp, nil
}

func MustNewInterpreter() *Interpreter {
	interp, err := NewInterpreter()
	if err != nil {
		panic(err)
	}
	return interp
}

func (interp *Interpreter) LexAndEvaluate(str string) (term.Term, error) {
	t, err := lexer.Lex(str)
	if err != nil {
		return t, err
	}
	err = interp.Evaluate(&t)
	return t, err
}

func (interp *Interpreter) MustLexAndEvaluate(str string) term.Term {
	t, err := interp.LexAndEvaluate(str)
	if err != nil {
		panic(err)
	}
	return t
}

func (interp *Interpreter) Evaluate(t *term.Term) error {
	if t == nil || t.IsEmpty() {
		return errors.New("Empty term")
	}
	// term has just one child, and no operation
	if len(t.Children) == 1 && len(t.Text) == 0 {
		*t = t.Children[0]
		return interp.Evaluate(t)
	}

	if len(t.Text) > 0 {
		interp.Logger.Printf("Evaluating `%s'\n", t.String())
		fnName := strings.ToLower(t.Text)
		if fn, ok := interp.functionMap[fnName]; ok {

			for i := range t.Children {
				if err := interp.Evaluate(&t.Children[i]); err != nil {
					return errors.Errorf("Error in child %s: %v", t.Children[i].Text, err)
				}
			}
			interp.Logger.Printf("Running fn `%s'\n", fnName)
			result, err := fn(&interp.Interpreter, t.Children...)
			if err != nil {
				return errors.Errorf("Error in function %s: %v", fnName, err)
			}
			interp.Logger.Printf("Updating value to `%s'\n", result)
			t.Update(result)
		}
	}

	return nil
}
