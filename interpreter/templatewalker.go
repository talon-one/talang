package interpreter

import (
	"strings"

	"github.com/talon-one/talang/token"
)

type templateWalker struct {
	interp *Interpreter
	pos    int
}

func (f *templateWalker) Next() *TaTemplate {
n:
	if f.pos >= len(f.interp.Templates) {
		if f.interp.Parent != nil {
			f.pos = 0
			f.interp = f.interp.Parent
			goto n
		}
		return nil
	}
	fn := &f.interp.Templates[f.pos]
	f.pos++
	return fn
}

type templateToRunWalker struct {
	templateWalker templateWalker
	Interpreter    *Interpreter
	Token          *token.TaToken
	Level          int
	lowerFuncName  string
	atomChildren   map[int][]token.Kind
	argumentCount  int
}

func newtemplateToRunWalker(interpreter *Interpreter, token *token.TaToken, level int) *templateToRunWalker {
	return &templateToRunWalker{
		templateWalker: templateWalker{
			interp: interpreter,
		},
		Interpreter:   interpreter,
		Level:         level,
		Token:         token,
		lowerFuncName: strings.ToLower(token.String),
		argumentCount: len(token.Children),
	}
}

func (f *templateToRunWalker) getChild(index int) []token.Kind {
	if f.atomChildren == nil {
		f.atomChildren = make(map[int][]token.Kind)
	}
	if child, ok := f.atomChildren[index]; ok {
		return child
	}

	walker := newtemplateToRunWalker(f.Interpreter, f.Token.Children[index], f.Level+1)

	var fns []token.Kind

	for fn := walker.Next(); fn != nil; fn = walker.Next() {
		fns = append(fns, fn.Returns)
	}

	f.atomChildren[index] = fns

	return fns
}

func (f *templateToRunWalker) Next() *TaTemplate {
outerloop:
	for fn := f.templateWalker.Next(); fn != nil; fn = f.templateWalker.Next() {
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
