package interpreter

import (
	"fmt"
	"go/ast"
	"strings"

	"reflect"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"

	"github.com/talon-one/talang/interpreter/shared"
	lexer "github.com/talon-one/talang/lexer"
)

type Interpreter struct {
	shared.Interpreter
}

func NewInterpreter() (*Interpreter, error) {
	interp := Interpreter{}
	if err := interp.registerCoreFunctions(); err != nil {
		return nil, err
	}
	interp.Binding = make(map[string]shared.Binding)
	return &interp, nil
}

func MustNewInterpreter() *Interpreter {
	interp, err := NewInterpreter()
	if err != nil {
		panic(err)
	}
	return interp
}

func (interp *Interpreter) LexAndEvaluate(str string) (*block.Block, error) {
	t, err := lexer.Lex(str)
	if err != nil {
		return t, err
	}
	err = interp.Evaluate(t)
	return t, err
}

func (interp *Interpreter) MustLexAndEvaluate(str string) *block.Block {
	t, err := interp.LexAndEvaluate(str)
	if err != nil {
		panic(err)
	}
	return t
}

func (interp *Interpreter) Evaluate(b *block.Block) error {
	if b == nil || b.IsEmpty() {
		return errors.New("Empty term")
	}

	childCount := len(b.Children)

	// term has just one child, and no operation
	if childCount == 1 && len(b.String) == 0 {
		*b = *b.Children[0]
		return interp.Evaluate(b)
	}

	if len(b.String) > 0 {
		if interp.Logger != nil {
			interp.Logger.Printf("Evaluating `%s'\n", b.Stringify())
		}
		stopProcessing, err := interp.callFunc(b)
		if err != nil {
			return err
		}
		if stopProcessing {
			return nil
		}
	}

	if interp.Parent != nil {
		sharedInterp := Interpreter{
			Interpreter: *interp.Parent,
		}
		return sharedInterp.Evaluate(b)
	}
	return nil
}

func (interp *Interpreter) MustEvaluate(b *block.Block) {
	if err := interp.Evaluate(b); err != nil {
		panic(err)
	}
}

type notMatchingDetail int

const (
	invalidName               notMatchingDetail = iota
	invalidSignature          notMatchingDetail = iota
	errorInChildrenEvaluation notMatchingDetail = iota
)

func (interp *Interpreter) matchesSignature(sig *shared.CommonSignature, lowerName string, args []*block.Block) (bool, notMatchingDetail, []*block.Block, error) {
	if sig.Name != lowerName {
		return false, invalidName, nil, nil
	}
	// make a copy of the children
	children := make([]*block.Block, len(args))
	for i, child := range args {
		children[i] = new(block.Block)
		*children[i] = *child
	}

	// evaluate children if needed to
	i := 0
	for ; i < len(sig.Arguments) && i < len(children); i++ {
		if sig.Arguments[i]&block.BlockKind == 0 && children[i].IsBlock() {
			if err := interp.Evaluate(children[i]); err != nil {
				return false, errorInChildrenEvaluation, nil, errors.Errorf("Error in child %s: %v", children[i].String, err)
			}
		}
	}
	if sig.IsVariadic {
		lastArgumentIndex := len(sig.Arguments) - 1
		// evaluate the rest
		for ; i < len(children); i++ {
			if sig.Arguments[lastArgumentIndex]&block.BlockKind == 0 && children[i].IsBlock() {
				if err := interp.Evaluate(children[i]); err != nil {
					return false, errorInChildrenEvaluation, nil, errors.Errorf("Error in child %s: %v", children[i].String, err)
				}
			}
		}
	}
	if sig.MatchesArguments(block.Arguments(children)) {
		return true, 0, children, nil
	}
	return false, invalidSignature, nil, nil
}

func (interp *Interpreter) callFunc(b *block.Block) (bool, error) {
	functions := interp.AllFunctions()
	blockText := strings.ToLower(b.String)
	// iterate trough all functions
	for n := 0; n < len(functions); n++ {
		fn := functions[n]
		run, notMatchingDetail, children, err := interp.matchesSignature(&fn.CommonSignature, blockText, b.Children)

		if !run {
			if interp.Logger != nil {
				switch notMatchingDetail {
				case invalidSignature:
					interp.Logger.Printf("NOT Running function `%s' (not matching signature)\n", fn.String())
				case errorInChildrenEvaluation:
					interp.Logger.Printf("NOT Running function `%s' (errors in child evaluation)\n", fn.String())
				}
			}
			if err != nil {
				return false, err
			}
			continue
		}
		// paranoid check
		if err != nil {
			return false, err
		}
		if interp.Logger != nil {
			interp.Logger.Printf("Running function `%s' with `%v'\n", fn.String(), block.BlockArguments(children).ToHumanReadable())
		}
		result, err := fn.Func(&interp.Interpreter, children...)
		if err != nil {
			if interp.Logger != nil {
				interp.Logger.Printf("Error in function %s: %v", fn.Name, err)
			}
			return false, errors.Errorf("Error in function %s: %v", fn.Name, err)
		}
		if result == nil {
			result = block.NewNull()
		}
		if fn.CommonSignature.Returns&result.Kind != result.Kind {
			if interp.Logger != nil {
				interp.Logger.Printf("Unexpected return type for %s: was `%s' expected %s", fn.Name, result.Kind.String(), fn.CommonSignature.Returns.String())
			}
			return false, errors.Errorf("Unexpected return type for %s: was `%s' expected %s", fn.Name, result.Kind.String(), fn.CommonSignature.Returns.String())
		}
		if interp.Logger != nil {
			interp.Logger.Printf("Updating value to `%s' (%s)\n", result.String, result.Kind.String())
		}
		b.Update(result)
		if b.IsBlock() {
			return true, interp.Evaluate(b)
		}
		return true, nil
	}
	return false, nil
}

func (interp *Interpreter) Set(key string, value shared.Binding) {
	interp.Binding[key] = value
}

func genericSetConv(value interface{}) (*shared.Binding, error) {
	reflectValue := reflect.ValueOf(value)
	reflectType := reflectValue.Type()
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflectValue.Elem()
	}

	switch reflectType.Kind() {
	case reflect.Struct:
		var bind shared.Binding
		bind.Children = make(map[string]shared.Binding)
		for i := 0; i < reflectType.NumField(); i++ {
			if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
				structValue, err := genericSetConv(reflectValue.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				bind.Children[fieldStruct.Name] = *structValue
			}
		}
		return &bind, nil
	case reflect.Int:
		fallthrough
	case reflect.Int8:
		fallthrough
	case reflect.Int16:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		fallthrough
	case reflect.Uint:
		fallthrough
	case reflect.Uint8:
		fallthrough
	case reflect.Uint16:
		fallthrough
	case reflect.Uint32:
		fallthrough
	case reflect.Uint64:
		return &shared.Binding{
			Value: block.New(fmt.Sprintf("%v", value)),
		}, nil
	case reflect.String:
		return &shared.Binding{
			Value: block.NewString(value.(string)),
		}, nil
	case reflect.Bool:
		return &shared.Binding{
			Value: block.NewBool(value.(bool)),
		}, nil
	}
	return nil, errors.Errorf("Unknown type `%T'", value)
}

func (interp *Interpreter) GenericSet(key string, value interface{}) error {
	binding, err := genericSetConv(value)
	if err != nil {
		return err
	}

	interp.Binding[key] = *binding
	return nil
}

func (interp *Interpreter) NewScope() *Interpreter {
	i := Interpreter{}
	i.Binding = make(map[string]shared.Binding)
	i.Parent = &interp.Interpreter
	i.Logger = interp.Logger
	// we need to register binding and template on this scope, because it uses its own scopes
	i.Functions = []shared.TaFunction{templateSignature(&i), bindingSignature}
	return &i
}
