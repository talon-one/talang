package interpreter

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/interpreter/corefn/cmp"
	"github.com/talon-one/talang/interpreter/corefn/math"
	"github.com/talon-one/talang/interpreter/corefn/misc"
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

func (interp *Interpreter) registerCoreFunctions() error {
	// simple mathematics
	interp.functions = append(interp.functions, math.AllOperations()...)
	// compare functions
	interp.functions = append(interp.functions, cmp.AllOperations()...)

	interp.functions = append(interp.functions, misc.Misc)
	return nil
}

func (interp *Interpreter) RemoveAllFunctions() error {
	interp.functions = []shared.TaSignature{}
	return nil
}
