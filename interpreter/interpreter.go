package interpreter

import (
	"context"
	"log"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/token"

	lexer "github.com/talon-one/talang/lexer"
)

type Interpreter struct {
	Binding           *token.TaToken
	Context           context.Context
	Parent            *Interpreter
	Functions         []TaFunction
	Templates         []TaTemplate
	Logger            *log.Logger
	IsDryRun          bool
	MaxRecursiveLevel *int
}

func NewInterpreter() (*Interpreter, error) {
	var interp Interpreter
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

func (interp *Interpreter) LexAndEvaluate(str string) (*token.TaToken, error) {
	t, err := lexer.Lex(str)
	if err != nil {
		return t, err
	}
	err = interp.Evaluate(t)
	return t, err
}

func (interp *Interpreter) MustLexAndEvaluate(str string) *token.TaToken {
	t, err := interp.LexAndEvaluate(str)
	if err != nil {
		panic(err)
	}
	return t
}

func (interp *Interpreter) Evaluate(b *token.TaToken) error {
	return interp.evaluate(b, 0)
}
func (interp *Interpreter) evaluate(b *token.TaToken, level int) error {
	if interp.MaxRecursiveLevel != nil && level > *interp.MaxRecursiveLevel {
		return &MaxRecursiveLevelReachedError{*interp.MaxRecursiveLevel}
	}
	if b == nil || b.IsEmpty() {
		return errors.New("Empty term")
	}

	if len(b.String) > 0 {
		var oldPrefix string
		if interp.Logger != nil {
			interp.Logger.Printf("Evaluating `%s'\n", b.Stringify())
			oldPrefix = interp.Logger.Prefix()
			interp.Logger.SetPrefix(oldPrefix + ">")
		}
		stopProcessing, err := interp.callFunc(b, level+1)
		if err != nil {
			return err
		}
		if interp.Logger != nil {
			interp.Logger.SetPrefix(oldPrefix)
		}
		if stopProcessing {
			return nil
		}
	} else if b.IsBlock() {
		size := len(b.Children)
		for i := 0; i < size; i++ {
			if err := interp.evaluate(b.Children[i], level+1); err != nil {
				return err
			}
		}
		token.Copy(b, b.Children[size-1])
		return nil
	}

	if interp.Parent != nil {
		return interp.Parent.evaluate(b, level+1)
	}
	return nil
}

func (interp *Interpreter) MustEvaluate(b *token.TaToken) {
	if err := interp.Evaluate(b); err != nil {
		panic(err)
	}
}

func (interp *Interpreter) callFunc(b *token.TaToken, level int) (bool, error) {
	if !b.IsBlock() {
		return false, nil
	}

	var collectedErrors []error
	walker := newFuncToRunWalker(interp, b, level+1)

nextfunc:
	for fn := walker.Next(); fn != nil; fn = walker.Next() {
		var result *token.TaToken
		if interp.IsDryRun {
			if interp.Logger != nil {
				interp.Logger.Printf("DryRunning function `%s' with `%v'\n", fn.String(), token.TokenArguments(b.Children).ToHumanReadable())
			}
			result = &token.TaToken{
				Kind: fn.Returns,
			}
		} else {
			// evaluating children
			var children []*token.TaToken

			if hasTokenBlock(b) {
				// make a copy of the children
				children = make([]*token.TaToken, len(b.Children))
				for i, child := range b.Children {
					children[i] = new(token.TaToken)
					token.Copy(children[i], child)
				}
			} else {
				children = b.Children
			}

			sigargc := len(fn.CommonSignature.Arguments)

			for i, j := 0, 0; i < len(children); i++ {
				if fn.CommonSignature.Arguments[j]&token.Token == 0 && children[i].IsBlock() {
					if err := interp.evaluate(children[i], level+1); err != nil {
						// children got an error
						return false, FunctionError{
							error:    err,
							function: fn,
						}
					}
				}
				j++
				if j >= sigargc {
					j = sigargc - 1
				}
			}

			// the children do not match after evaluation => goto next function
			if !fn.CommonSignature.MatchesArguments(token.Arguments(children)) {

				collectedErrors = append(collectedErrors, FunctionNotRanError{
					Reason:   errors.New("Not matching signature - after evaluation"),
					function: fn,
				})
				continue nextfunc
			}

			if interp.Logger != nil {
				interp.Logger.Printf("Running function `%s' with `%v'\n", fn.String(), token.TokenArguments(children).ToHumanReadable())
			}
			var err error
			result, err = fn.Func(interp, children...)
			// error in function
			if err != nil {
				return false, FunctionError{error: err, function: fn}
			}
			if result == nil {
				result = token.NewNull()
			}
		}

		// we ran the function and everything is okay
		// make sure the result matches our expected data
		if fn.CommonSignature.Returns&result.Kind != result.Kind {
			err := errors.Errorf("Unexpected return type for `%s': was `%s' expected `%s'", fn.Name, result.Kind.String(), fn.CommonSignature.Returns.String())
			if interp.Logger != nil {
				interp.Logger.Println(err)
			}
			// if not => goto next function
			collectedErrors = append(collectedErrors, err)
			err = nil
			continue
		}
		if interp.Logger != nil {
			interp.Logger.Printf("Updating value to `%s' (%s)\n", result.Stringify(), result.Kind.String())
		}
		token.Copy(b, result)
		if b.IsBlock() {
			return true, interp.evaluate(b, level+1)
		}
		return true, nil
	}
	// we found no matching function OR all functions failed
	err := FunctionNotFoundError{CollectedErrors: collectedErrors, token: b}
	if interp.Logger != nil {
		interp.Logger.Println(err)
	}
	return false, err
}

func (interp *Interpreter) Get(key string) *token.TaToken {
	if interp.Binding != nil {
		return interp.Binding.MapItem(key)
	}
	return token.NewNull()
}

func (interp *Interpreter) Set(key string, value *token.TaToken) {
	if interp.Binding == nil {
		interp.Binding = token.NewMap(map[string]*token.TaToken{})
	}
	interp.Binding.SetMapItem(key, value)
}

func (interp *Interpreter) NewScope() *Interpreter {
	i := Interpreter{}
	i.Parent = interp
	i.Logger = interp.Logger
	i.MaxRecursiveLevel = interp.MaxRecursiveLevel
	// we need to register binding and template on this scope, because it uses its own scopes
	i.Functions = []TaFunction{templateSignature, setTemplateSignature, bindingSignature, setBindingSignature}
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

func hasTokenBlock(tkn *token.TaToken) bool {
	for _, child := range tkn.Children {
		if child.IsBlock() {
			return true
		}
	}
	return false
}
