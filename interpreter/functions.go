package interpreter

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/interpreter/corefn/cmp"
	"github.com/talon-one/talang/interpreter/corefn/math"
	"github.com/talon-one/talang/interpreter/internal"
	"github.com/talon-one/talang/term"
)

type TaFunc func(interp *internal.Interpreter, args ...term.Term) (string, error)

func (interp *Interpreter) RegisterFunction(name string, fn TaFunc) error {
	name = strings.ToLower(name)
	if _, ok := interp.functionMap[name]; ok {
		return errors.Errorf("Function `%s' is already registered", name)
	}
	interp.functionMap[name] = fn
	return nil
}

func (interp *Interpreter) UpdateFunction(name string, fn TaFunc) error {
	name = strings.ToLower(name)
	if _, ok := interp.functionMap[name]; !ok {
		return errors.Errorf("Function `%s' is not registered", name)
	}
	interp.functionMap[name] = fn
	return nil
}

func (interp *Interpreter) RemoveFunction(name string) error {
	name = strings.ToLower(name)
	if _, ok := interp.functionMap[name]; !ok {
		return errors.Errorf("Function `%s' is not registered", name)
	}
	delete(interp.functionMap, name)
	return nil
}

func (interp *Interpreter) registerCoreFunctions() error {

	fns := []struct {
		name string
		fn   TaFunc
	}{
		// simple mathematics
		{"+", math.Add},
		{"-", math.Sub},
		{"*", math.Mul},
		{"/", math.Div},
		{"mod", math.Mod},
		{"ceil", math.Ceil},
		{"floor", math.Floor},

		// compare
		{"=", cmp.Equal},
		{"!=", cmp.NotEqual},
		{">", cmp.GreaterThan},
		{">=", cmp.GreaterThanOrEqual},
		{"<", cmp.LessThan},
		{"<=", cmp.LessThanOrEqual},
		{"between", cmp.Between},
	}

	for i := 0; i < len(fns); i++ {
		interp.functionMap[strings.ToLower(fns[i].name)] = fns[i].fn
	}
	return nil
}
