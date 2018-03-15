package interpreter

import (
	"context"
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"

	lexer "github.com/talon-one/talang/lexer"
)

type Interpreter struct {
	Binding   *block.Block
	Context   context.Context
	Parent    *Interpreter
	Functions []TaFunction
	Templates []TaTemplate
	Logger    *log.Logger
}

func NewInterpreter() (*Interpreter, error) {
	interp := Interpreter{}
	if err := interp.registerCoreFunctions(); err != nil {
		return nil, err
	}
	if err := interp.registerCoreTemplates(); err != nil {
		return nil, err
	}
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
		var oldPrefix string
		if interp.Logger != nil {
			interp.Logger.Printf("Evaluating `%s'\n", b.Stringify())
			oldPrefix = interp.Logger.Prefix()
			interp.Logger.SetPrefix(oldPrefix + ">")
		}
		stopProcessing, err := interp.callFunc(b)
		if err != nil {
			return err
		}
		if interp.Logger != nil {
			interp.Logger.SetPrefix(oldPrefix)
		}
		if stopProcessing {
			return nil
		}
	}

	if interp.Parent != nil {
		return interp.Parent.Evaluate(b)
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

func (interp *Interpreter) matchesSignature(sig *CommonSignature, lowerName string, args []*block.Block) (bool, notMatchingDetail, []*block.Block, error) {
	if sig.lowerName != lowerName {
		return false, invalidName, nil, nil
	}

	var children []*block.Block
	argc := len(args)
	sigargc := len(sig.Arguments)
	if !sig.IsVariadic {
		if sigargc != argc {
			return false, invalidSignature, nil, nil
		}
	} else {
		if sigargc-1 > argc {
			return false, invalidSignature, nil, nil
		}
	}

	willEvaluate := false
	for i, j := 0, 0; i < argc; i++ {
		if sig.Arguments[j]&block.BlockKind == 0 && args[i].IsBlock() {
			willEvaluate = true
			break
		}
		j++
		if j >= sigargc {
			j = sigargc - 1
		}
	}

	if willEvaluate {
		// make a copy of the children
		children = make([]*block.Block, argc)
		for i, child := range args {
			children[i] = new(block.Block)
			*children[i] = *child
		}
	} else {
		children = args
	}

	for i, j := 0, 0; i < len(children); i++ {
		if sig.Arguments[j]&block.BlockKind == 0 && children[i].IsBlock() {
			if err := interp.Evaluate(children[i]); err != nil {
				return false, errorInChildrenEvaluation, nil, errors.Errorf("Error in child %s: %v", children[i].String, err)
			}
		}
		j++
		if j >= sigargc {
			j = sigargc - 1
		}
	}

	if sig.MatchesArguments(block.Arguments(children)) {
		return true, 0, children, nil
	}
	return false, invalidSignature, nil, nil
}

func (interp *Interpreter) callFunc(b *block.Block) (bool, error) {
	walker := funcWalker{interp: interp}
	blockText := strings.ToLower(b.String)
	// iterate trough all functions
	for fn := walker.Next(); fn != nil; fn = walker.Next() {
		run, detail, children, err := interp.matchesSignature(&fn.CommonSignature, blockText, b.Children)

		if !run {
			if interp.Logger != nil {
				switch detail {
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
		result, err := fn.Func(interp, children...)
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
			interp.Logger.Printf("Updating value to `%s' (%s)\n", result.Stringify(), result.Kind.String())
		}
		block.Copy(b, result)
		if b.IsBlock() {
			return true, interp.Evaluate(b)
		}
		return true, nil
	}
	if interp.Logger != nil {
		interp.Logger.Printf("Found no func `%s' in interpreter instance\n", blockText)
	}
	return false, nil
}

func (interp *Interpreter) Get(key string) *block.Block {
	if interp.Binding != nil {
		return interp.Binding.MapItem(key)
	}
	return block.NewNull()
}

func (interp *Interpreter) Set(key string, value *block.Block) {
	if interp.Binding == nil {
		interp.Binding = block.NewMap(map[string]*block.Block{})
	}
	interp.Binding.SetMapItem(key, value)
}

func genericSetConv(value interface{}) (*block.Block, error) {
	reflectValue := reflect.ValueOf(value)
	reflectType := reflectValue.Type()
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectValue = reflectValue.Elem()
	}

	switch reflectType.Kind() {
	case reflect.Struct:
		m := make(map[string]*block.Block)
		for i := 0; i < reflectType.NumField(); i++ {
			if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) {
				structValue, err := genericSetConv(reflectValue.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				m[fieldStruct.Name] = structValue
			}
		}
		return block.NewMap(m), nil
	case reflect.Map:
		m := make(map[string]*block.Block)
		if reflectType.Key().Kind() != reflect.String {
			return nil, errors.New("A different key than `string' is not supported")
		}
		for _, key := range reflectValue.MapKeys() {
			var err error
			m[key.String()], err = genericSetConv(reflectValue.MapIndex(key).Interface())
			if err != nil {
				return nil, err
			}
		}
		return block.NewMap(m), nil
	case reflect.Slice:
		size := reflectValue.Len()
		s := make([]*block.Block, size, size)
		for i := 0; i < size; i++ {
			var err error
			s[i], err = genericSetConv(reflectValue.Index(i).Interface())
			if err != nil {
				return nil, err
			}
		}
		return block.NewList(s...), nil
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
		return block.NewDecimalFromString(fmt.Sprintf("%v", value)), nil
	case reflect.String:
		return block.NewString(value.(string)), nil
	case reflect.Bool:
		return block.NewBool(value.(bool)), nil
	case reflect.Float32:
		return block.NewDecimalFromFloat(float64(value.(float32))), nil
	case reflect.Float64:
		return block.NewDecimalFromFloat(value.(float64)), nil
	}
	return nil, errors.Errorf("Unknown type `%T'", value)
}

func (interp *Interpreter) GenericSet(key string, value interface{}) error {
	block, err := genericSetConv(value)
	if err != nil {
		return err
	}

	if len(key) == 0 {
		interp.Binding = block
	} else {
		interp.Set(key, block)
	}
	return nil
}

func (interp *Interpreter) NewScope() *Interpreter {
	i := Interpreter{}
	i.Parent = interp
	i.Logger = interp.Logger
	// we need to register binding and template on this scope, because it uses its own scopes
	i.Functions = []TaFunction{templateSignature, bindingSignature, setBindingSignature}
	return &i
}

func (interp *Interpreter) AllFunctions() (functions []TaFunction) {
	if len(interp.Functions) > 0 {
		functions = append(functions, interp.Functions...)
	}
	if interp.Parent != nil {
		functions = append(functions, interp.Parent.AllFunctions()...)
	}
	return functions
}

func (interp *Interpreter) AllTemplates() (templates []TaTemplate) {
	if len(interp.Templates) > 0 {
		templates = append(templates, interp.Templates...)
	}
	if interp.Parent != nil {
		templates = append(templates, interp.Parent.AllTemplates()...)
	}
	return templates
}

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
