package shared

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/talon-one/talang/block"
)

type Binding struct {
	Value    *block.Block
	Children map[string]Binding
}

type Interpreter struct {
	Binding   map[string]Binding
	Context   context.Context
	Parent    *Interpreter
	Functions []TaFunction
	Templates []TaTemplate
	Logger    *log.Logger
}

type TaFunc func(*Interpreter, ...*block.Block) (*block.Block, error)

type CommonSignature struct {
	IsVariadic  bool
	Arguments   []block.Kind
	Name        string
	Returns     block.Kind
	Description string
}

type TaFunction struct {
	CommonSignature
	Func TaFunc `json:"-"`
}

type TaTemplate struct {
	CommonSignature
	Template block.Block `json:"-"`
}

func (s *CommonSignature) String() string {
	var args string
	if argc := len(s.Arguments); argc > 0 {
		argv := make([]string, argc)
		for i := 0; i < argc; i++ {
			argv[i] = s.Arguments[i].String()
		}
		args = strings.Join(argv, ", ")
	}

	var variadic string
	if s.IsVariadic {
		variadic = "..."
	}
	return fmt.Sprintf("%s(%s)%s", s.Name, args, variadic)
}

func (a *CommonSignature) Equal(b *CommonSignature) bool {
	if a.IsVariadic != b.IsVariadic {
		return false
	}
	if len(a.Arguments) != len(b.Arguments) {
		return false
	}
	for i, arg := range a.Arguments {
		if b.Arguments[i] != arg {
			return false
		}
	}
	return a.Name == b.Name
}

func (sig *CommonSignature) MatchesArguments(args []block.Kind) bool {
	if !sig.IsVariadic {
		if len(args) != len(sig.Arguments) {
			return false
		}
		for i, kind := range args {
			if sig.Arguments[i]&kind == 0 {
				return false
			}
		}
		return true
	}
	sigArgc := len(sig.Arguments) - 1
	for i, j := 0, 0; i < len(args); i++ {
		if sig.Arguments[j]&args[i] == 0 {
			return false
		}
		if i < sigArgc {
			j++
		}
	}
	return true
}

func (s *TaFunction) String() string {
	return s.CommonSignature.String()
}

func (a *TaFunction) Equal(b *TaFunction) bool {
	return a.CommonSignature.Equal(&b.CommonSignature)
}

func (s *TaFunction) MatchesArguments(args []block.Kind) bool {
	return s.CommonSignature.MatchesArguments(args)
}

func (s *TaTemplate) String() string {
	return s.CommonSignature.String()
}

func (a *TaTemplate) Equal(b *TaTemplate) bool {
	return a.CommonSignature.Equal(&b.CommonSignature)
}

func (s *TaTemplate) MatchesArguments(args []block.Kind) bool {
	return s.CommonSignature.MatchesArguments(args)
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
