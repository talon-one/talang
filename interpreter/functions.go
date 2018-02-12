package interpreter

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/corefn/cmp"
	"github.com/talon-one/talang/interpreter/corefn/math"
	"github.com/talon-one/talang/interpreter/corefn/misc"
	stringpkg "github.com/talon-one/talang/interpreter/corefn/string"
	"github.com/talon-one/talang/interpreter/shared"
)

func (interp *Interpreter) RegisterFunction(signature shared.TaSignature) error {
	signature.Name = strings.ToLower(signature.Name)
	if interp.GetFunction(signature) != nil {
		return errors.Errorf("Function `%s' is already registered", signature.Name)
	}
	interp.functions = append(interp.functions, signature)
	return nil
}

func (interp *Interpreter) UpdateFunction(signature shared.TaSignature) error {
	signature.Name = strings.ToLower(signature.Name)
	if s := interp.GetFunction(signature); s != nil {
		*s = signature
		return nil
	}
	return errors.Errorf("Function `%s' is not registered", signature.Name)
}

func (interp *Interpreter) RemoveFunction(signature shared.TaSignature) error {
	signature.Name = strings.ToLower(signature.Name)
	for i := 0; i < len(interp.functions); i++ {
		if interp.functions[i].Equal(&signature) {
			fns := interp.functions[:i]
			interp.functions = append(fns, interp.functions[i+1:]...)
			return nil
		}
	}
	return errors.Errorf("Function `%s' is not registered", signature.Name)
}

func (interp *Interpreter) GetFunction(signature shared.TaSignature) *shared.TaSignature {
	signature.Name = strings.ToLower(signature.Name)
	for i := 0; i < len(interp.functions); i++ {
		if interp.functions[i].Equal(&signature) {
			return &interp.functions[i]
		}
	}
	return nil
}

func (interp *Interpreter) Functions() []shared.TaSignature {
	fns := make([]shared.TaSignature, len(interp.functions))
	for i, fn := range interp.functions {
		fns[i] = fn
	}
	return fns
}

func (interp *Interpreter) registerCoreFunctions() error {
	// simple mathematics
	interp.functions = append(interp.functions, math.AllOperations()...)
	// compare functions
	interp.functions = append(interp.functions, cmp.AllOperations()...)

	// string functions
	interp.functions = append(interp.functions, stringpkg.AllOperations()...)

	interp.functions = append(interp.functions, misc.Misc)

	// binding
	interp.functions = append(interp.functions, shared.TaSignature{
		Name:       ".",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AtomKind,
		},
		Returns: block.BlockKind,
		Func: func(interp *shared.Interpreter, args []*block.Block) (*block.Block, error) {
			bindMap := interp.Binding
			var value *block.Block
			for i := 0; i < len(args); i++ {
				if binding, ok := bindMap[args[i].Text]; ok {
					bindMap = binding.Children
					value = binding.Value
				} else {
					// join args
					qualifiers := make([]string, len(args))
					for j, arg := range args {
						qualifiers[j] = arg.Text
					}
					return nil, errors.Errorf("Unable to find `%s'", strings.Join(qualifiers, "."))
				}
			}
			//
			return value, nil
		},
	})
	return nil
}

func (interp *Interpreter) RemoveAllFunctions() error {
	interp.functions = []shared.TaSignature{}
	return nil
}
