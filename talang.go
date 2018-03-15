//go:generate go run generate_corefn.go -pkg=talang
package talang

// for convenience all important stuff is here

import (
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/lexer"
)

type Interpreter struct {
	interpreter.Interpreter
}

func Lex(str string) (*block.TaToken, error) {
	return lexer.Lex(str)
}

func MustLex(str string) *block.TaToken {
	block, err := lexer.Lex(str)
	if err != nil {
		panic(err)
	}
	return block
}

func Parse(str string) (*block.TaToken, error) {
	return Lex(str)
}
func MustParse(str string) *block.TaToken {
	return MustLex(str)
}

func NewInterpreter() (*Interpreter, error) {
	interp, err := interpreter.NewInterpreter()
	return &Interpreter{*interp}, err
}

func MustNewInterpreter() *Interpreter {
	interp := interpreter.MustNewInterpreter()
	return &Interpreter{*interp}
}

func (interp *Interpreter) LexAndEvaluate(str string) (*block.TaToken, error) {
	return interp.Interpreter.LexAndEvaluate(str)
}

func (interp *Interpreter) MustLexAndEvaluate(str string) *block.TaToken {
	return interp.Interpreter.MustLexAndEvaluate(str)
}

func (interp *Interpreter) Evaluate(b *block.TaToken) error {
	return interp.Interpreter.Evaluate(b)
}

func (interp *Interpreter) MustEvaluate(b *block.TaToken) {
	interp.Interpreter.MustEvaluate(b)
}
