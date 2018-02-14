package talang

// for convenience all important stuff is here

import (
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/lexer"
)

type Interpreter interpreter.Interpreter

func Lex(str string) (*block.Block, error) {
	return lexer.Lex(str)
}

func MustLex(str string) *block.Block {
	block, err := lexer.Lex(str)
	if err != nil {
		panic(err)
	}
	return block
}

func Parse(str string) (*block.Block, error) {
	return Lex(str)
}
func MustParse(str string) *block.Block {
	return MustLex(str)
}

func NewInterpreter() (*Interpreter, error) {
	interp, err := interpreter.NewInterpreter()
	var i Interpreter
	i = Interpreter(*interp)
	return &i, err
}
