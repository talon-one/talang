package interpreter

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/pkg/errors"
	"github.com/talon-one/talang/token"

	lexer "github.com/talon-one/talang/lexer"
)

type EvaluationMode int

const (
	Unsafe EvaluationMode = iota
	Safe   EvaluationMode = iota
)

type Interpreter struct {
	Binding           *token.TaToken
	Context           context.Context
	Parent            *Interpreter
	Functions         []TaFunction
	Templates         []TaTemplate
	Logger            *log.Logger
	IsDryRun          bool
	EvaluationMode    EvaluationMode
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
		return errors.Errorf("Max Recursive level (%d) reached", *interp.MaxRecursiveLevel)
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

type notMatchingDetail int

const (
	invalidName               notMatchingDetail = iota
	invalidSignature          notMatchingDetail = iota
	errorInChildrenEvaluation notMatchingDetail = iota
)

func (interp *Interpreter) matchesSignature(sig *CommonSignature, lowerName string, args []*token.TaToken, level int) (bool, notMatchingDetail, []*token.TaToken, error) {
	if sig.lowerName != lowerName {
		return false, invalidName, nil, nil
	}

	var children []*token.TaToken
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
		if sig.Arguments[j]&token.Token == 0 && args[i].IsBlock() {
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
		children = make([]*token.TaToken, argc)
		for i, child := range args {
			children[i] = new(token.TaToken)
			*children[i] = *child
		}
	} else {
		children = args
	}

	for i, j := 0, 0; i < len(children); i++ {
		if sig.Arguments[j]&token.Token == 0 && children[i].IsBlock() {
			if err := interp.evaluate(children[i], level+1); err != nil {
				return false, errorInChildrenEvaluation, nil, errors.Errorf("Error in child %s: %v", children[i].String, err)
			}
		}
		j++
		if j >= sigargc {
			j = sigargc - 1
		}
	}

	if sig.MatchesArguments(token.Arguments(children)) {
		return true, 0, children, nil
	}
	return false, invalidSignature, nil, nil
}

func (interp *Interpreter) callFunc(b *token.TaToken, level int) (bool, error) {
	walker := funcWalker{interp: interp}
	blockText := strings.ToLower(b.String)
	var collectedErrors []error
	// iterate trough all functions
	for fn := walker.Next(); fn != nil; fn = walker.Next() {
		run, detail, children, err := interp.matchesSignature(&fn.CommonSignature, blockText, b.Children, level+1)

		if interp.EvaluationMode == Safe {
			switch detail {
			case invalidSignature:
				collectedErrors = append(collectedErrors, errors.Errorf(`  Expression %s doesn't match '%s'`, b.Stringify(), fn.String()))
			case errorInChildrenEvaluation:
				collectedErrors = append(collectedErrors, errors.Errorf("  NOT Running function `%s' (errors in child evaluation)", fn.String()))
			}
		}

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

		var result *token.TaToken
		if !interp.IsDryRun {
			if interp.Logger != nil {
				interp.Logger.Printf("Running function `%s' with `%v'\n", fn.String(), token.TokenArguments(children).ToHumanReadable())
			}
			result, err = fn.Func(interp, children...)
		} else {
			if interp.Logger != nil {
				interp.Logger.Printf("DryRunning function `%s' with `%v'\n", fn.String(), token.TokenArguments(children).ToHumanReadable())
			}
			result = &token.TaToken{
				Kind: fn.Returns,
			}
		}
		if err != nil {
			if interp.Logger != nil {
				interp.Logger.Printf("Error in function %s: %v", fn.Name, err)
			}
			return false, errors.Errorf("Error in function %s: %v", fn.Name, err)
		}
		if result == nil {
			result = token.NewNull()
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
		token.Copy(b, result)
		if b.IsBlock() {
			return true, interp.evaluate(b, level+1)
		}
		return true, nil
	}
	if interp.Logger != nil {
		interp.Logger.Printf("Found no func `%s' in interpreter instance\n", blockText)
	}
	if interp.EvaluationMode == Safe {
		var builder strings.Builder
		builder.WriteString(fmt.Sprintf("Found no eval function for %s", b.Stringify()))
		builder.WriteRune('\n')
		for i := 0; i < len(collectedErrors); i++ {
			builder.WriteString(collectedErrors[i].Error())
			builder.WriteRune('\n')
		}
		return false, errors.Errorf(builder.String())
	}
	return false, nil
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
