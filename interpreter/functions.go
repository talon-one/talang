package interpreter

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/interpreter/corefn/cmp"
	"github.com/talon-one/talang/interpreter/corefn/list"
	"github.com/talon-one/talang/interpreter/corefn/mapping"
	"github.com/talon-one/talang/interpreter/corefn/math"
	"github.com/talon-one/talang/interpreter/corefn/misc"
	stringpkg "github.com/talon-one/talang/interpreter/corefn/string"
	"github.com/talon-one/talang/interpreter/shared"
)

func (interp *Interpreter) RegisterFunction(signatures ...shared.TaFunction) error {
	for i := 0; i < len(signatures); i++ {
		signature := signatures[i]
		signature.Name = strings.ToLower(signature.Name)
		if interp.GetFunction(signature) != nil {
			return errors.Errorf("Function `%s' is already registered", signature.Name)
		}
		interp.Functions = append(interp.Functions, signature)
	}
	return nil
}

func (interp *Interpreter) UpdateFunction(signature shared.TaFunction) error {
	signature.Name = strings.ToLower(signature.Name)
	if s := interp.GetFunction(signature); s != nil {
		*s = signature
		return nil
	}
	return errors.Errorf("Function `%s' is not registered", signature.Name)
}

func (interp *Interpreter) RemoveFunction(signature shared.TaFunction) error {
	signature.Name = strings.ToLower(signature.Name)
	for i := 0; i < len(interp.Functions); i++ {
		if interp.Functions[i].Equal(&signature) {
			fns := interp.Functions[:i]
			interp.Functions = append(fns, interp.Functions[i+1:]...)
			return nil
		}
	}
	return errors.Errorf("Function `%s' is not registered", signature.Name)
}

func (interp *Interpreter) GetFunction(signature shared.TaFunction) *shared.TaFunction {
	signature.Name = strings.ToLower(signature.Name)
	for i := 0; i < len(interp.Functions); i++ {
		if interp.Functions[i].Equal(&signature) {
			return &interp.Functions[i]
		}
	}
	return nil
}

func (interp *Interpreter) registerCoreFunctions() error {
	// simple mathematics
	interp.Functions = append(interp.Functions, math.AllOperations()...)
	// compare functions
	interp.Functions = append(interp.Functions, cmp.AllOperations()...)

	// string functions
	interp.Functions = append(interp.Functions, stringpkg.AllOperations()...)

	// misc functions
	interp.Functions = append(interp.Functions, misc.AllOperations()...)

	// list functions
	interp.Functions = append(interp.Functions, list.AllOperations()...)

	// map functions
	interp.Functions = append(interp.Functions, mapping.AllOperations()...)

	// binding
	interp.Functions = append(interp.Functions, bindingSignature)

	// template
	interp.Functions = append(interp.Functions, templateSignature(interp))

	// sanitize name
	for i, f := range interp.Functions {
		interp.Functions[i].CommonSignature.Name = strings.ToLower(f.CommonSignature.Name)
	}
	return nil
}

func (interp *Interpreter) RemoveAllFunctions() error {
	interp.Functions = []shared.TaFunction{}
	return nil
}

var bindingSignature = shared.TaFunction{
	CommonSignature: shared.CommonSignature{
		Name:       ".",
		IsVariadic: true,
		Arguments: []block.Kind{
			block.AtomKind,
		},
		Returns:     block.AnyKind,
		Description: "Access a variable in the binding",
	},
	Func: bindingFunc,
}

func bindingFunc(interp *shared.Interpreter, args ...*block.Block) (*block.Block, error) {
	bindMap := interp.Binding
	var value *block.Block
	for i := 0; i < len(args); i++ {
		if binding, ok := bindMap[args[i].String]; ok {
			bindMap = binding.Children
			value = binding.Value
		} else {
			value = nil
		}
	}
	if value == nil {
		// lookup in parent
		if interp.Parent != nil {
			value, err := bindingFunc(interp.Parent, args...)
			if err == nil {
				return value, nil
			}
		}

		// join args
		qualifiers := make([]string, len(args))
		for j, arg := range args {
			qualifiers[j] = arg.String
		}
		return nil, errors.Errorf("Unable to find `%s'", strings.Join(qualifiers, "."))
	}
	return value, nil
}
