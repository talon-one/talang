package interpreter

import (
	"strings"

	"github.com/talon-one/talang/token"
)

type funcWalker struct {
	interp *Interpreter
	pos    int
}

func (f *funcWalker) Next() *TaFunction {
n:
	if f.pos >= len(f.interp.Functions) {
		if f.interp.Parent != nil {
			f.pos = 0
			f.interp = f.interp.Parent
			goto n
		}
		return nil
	}
	fn := &f.interp.Functions[f.pos]
	f.pos++
	return fn
}

type funcToRunWalker struct {
	funcWalker    funcWalker
	Interpreter   *Interpreter
	Token         *token.TaToken
	Level         int
	lowerFuncName string
	atomChildren  map[int][]token.Kind
	argumentCount int
}

func newFuncToRunWalker(interpreter *Interpreter, token *token.TaToken, level int) *funcToRunWalker {
	return &funcToRunWalker{
		funcWalker: funcWalker{
			interp: interpreter,
		},
		Interpreter:   interpreter,
		Level:         level,
		Token:         token,
		lowerFuncName: strings.ToLower(token.String),
		argumentCount: len(token.Children),
	}
}

func (f *funcToRunWalker) getChild(index int) []token.Kind {
	if f.atomChildren == nil {
		f.atomChildren = make(map[int][]token.Kind)
	}
	if child, ok := f.atomChildren[index]; ok {
		return child
	}

	walker := newFuncToRunWalker(f.Interpreter, f.Token.Children[index], f.Level+1)

	var fns []token.Kind

	for fn := walker.Next(); fn != nil; fn = walker.Next() {
		fns = append(fns, fn.Returns)
	}

	f.atomChildren[index] = fns

	return fns
}

func (f *funcToRunWalker) Next() *TaFunction {
outerloop:
	for fn := f.funcWalker.Next(); fn != nil; fn = f.funcWalker.Next() {
		if fn.lowerName != f.lowerFuncName {
			continue
		}

		fnArgc := len(fn.Arguments)
		if !fn.IsVariadic {
			if fnArgc != f.argumentCount {
				continue
			}
		} else {
			if fnArgc-1 > f.argumentCount {
				continue
			}
		}

		for i, j := 0, 0; i < f.argumentCount; i++ {
			child := f.Token.Children[i]
			childKinds := []token.Kind{child.Kind}
			// if the child is a block
			if child.IsBlock() {
				// but the fn does not accept a block
				if fn.Arguments[j]&token.Token == 0 {
					// use the result of the child
					childKinds = f.getChild(i)
				}
			}

			gotMatchingFunc := false
			for k := 0; k < len(childKinds); k++ {
				if fn.Arguments[j]&childKinds[k] != 0 {
					gotMatchingFunc = true
				}
			}

			if !gotMatchingFunc {
				continue outerloop
			}

			j++
			if j >= fnArgc {
				j = fnArgc - 1
			}
		}
		return fn
	}
	return nil
}
