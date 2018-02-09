package shared

import (
	"fmt"
	"strings"

	"github.com/talon-one/talang/block"
)

type Binding struct {
	Value    *block.Block
	Children map[string]Binding
}

type Interpreter struct {
	Binding map[string]Binding
}

type TaFunc func(*Interpreter, []*block.Block) (*block.Block, error)

type TaSignature struct {
	IsVariadic bool
	Arguments  []block.Kind
	Name       string
	Func       TaFunc
}

func (s *TaSignature) String() string {
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

func (a *TaSignature) Equal(b *TaSignature) bool {
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
	if a.Name != b.Name {
		return false
	}
	return true
}

func (sig *TaSignature) MatchesArguments(args []block.Kind) bool {
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
