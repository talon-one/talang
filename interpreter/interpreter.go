package interpreter

import (
	"fmt"
	"go/ast"
	"log"
	"strings"

	"reflect"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/block"

	"github.com/talon-one/talang/interpreter/shared"
	lexer "github.com/talon-one/talang/lexer"
)

type Interpreter struct {
	shared.Interpreter
	functions []shared.TaSignature
	Logger    *log.Logger
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
	if childCount == 1 && len(b.Text) == 0 {
		*b = *b.Children[0]
		return interp.Evaluate(b)
	}

	if len(b.Text) > 0 {
		if interp.Logger != nil {
			interp.Logger.Printf("Evaluating `%s'\n", b.String())
		}
		blockText := strings.ToLower(b.Text)
		// iterate trough all functions
		for n := 0; n < len(interp.functions); n++ {
			// if we have found a function that matches the name
			if interp.functions[n].Name == blockText {
				fn := interp.functions[n]

				// make a copy of the children

				children := make([]*block.Block, len(b.Children))
				for i, child := range b.Children {
					children[i] = new(block.Block)
					*children[i] = *child
				}

				// evaluate children if needed to
				i := 0
				for ; i < len(fn.Arguments) && i < len(children); i++ {
					if fn.Arguments[i]&block.AtomKind != 0 && children[i].IsBlock() {
						if err := interp.Evaluate(children[i]); err != nil {
							return errors.Errorf("Error in child %s: %v", children[i].Text, err)
						}
					}
				}
				if fn.IsVariadic {
					lastArgumentIndex := len(fn.Arguments) - 1
					// evaluate the rest
					for ; i < len(children); i++ {
						if fn.Arguments[lastArgumentIndex]&block.AtomKind != 0 && children[i].IsBlock() {
							if err := interp.Evaluate(children[i]); err != nil {
								return errors.Errorf("Error in child %s: %v", children[i].Text, err)
							}
						}
					}
				}
				if fn.MatchesArguments(block.Arguments(children)) {
					if interp.Logger != nil {
						humanReadableArguments := make([]string, len(children))
						for i, arg := range children {
							humanReadableArguments[i] = arg.String()
						}
						interp.Logger.Printf("Running fn `%s' with `%v'\n", fn.String(), strings.Join(humanReadableArguments, ", "))
					}
					result, err := fn.Func(&interp.Interpreter, children...)
					if err != nil {
						return errors.Errorf("Error in function %s: %v", fn.Name, err)
					}
					if result == nil {
						result = block.NewString("")
					}
					if interp.Logger != nil {
						interp.Logger.Printf("Updating value to `%s'\n", result)
					}
					b.Update(result)
					if b.IsBlock() {
						return interp.Evaluate(b)
					}
					break
				} else if interp.Logger != nil {
					interp.Logger.Printf("NOT Running fn `%s' (not matching signature)\n", fn.String())
				}
			}
		}

	}
	return nil
}

func (interp *Interpreter) MustEvaluate(b *block.Block) {
	if err := interp.Evaluate(b); err != nil {
		panic(err)
	}
}

func (interp *Interpreter) Set(key string, value shared.Binding) {
	interp.Binding[key] = value
}

func genericSetConv(value interface{}) (*shared.Binding, error) {
	reflectValue := reflect.ValueOf(value)
	reflectType := reflectValue.Type()
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
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
