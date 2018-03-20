package interpreter

import (
	"fmt"
	"strings"

	"github.com/talon-one/talang/token"
)

type TaFunc func(*Interpreter, ...*token.TaToken) (*token.TaToken, error)

type CommonSignature struct {
	IsVariadic  bool
	Arguments   []token.Kind
	Name        string
	lowerName   string
	Returns     token.Kind
	Description string
	Example     string
}

type TaFunction struct {
	CommonSignature
	Func TaFunc `json:"-"`
}

type TaTemplate struct {
	CommonSignature
	Template token.TaToken
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
	return fmt.Sprintf("%s(%s%s)", s.Name, args, variadic)
}

func NewCommonSignature(s string) *CommonSignature {
	size := len(s)
	if size <= 0 {
		return nil
	}
	bracketOpen := strings.IndexRune(s, '(')
	bracketClose := strings.IndexRune(s, ')')
	if bracketOpen <= 0 || bracketClose <= bracketOpen {
		return nil
	}

	var signature CommonSignature

	signature.Name = strings.TrimSpace(s[:bracketOpen])
	signature.lowerName = strings.ToLower(signature.Name)

	arguments := strings.Split(s[bracketOpen+1:bracketClose], ",")
	for i, l := 0, len(arguments)-1; i <= l; i++ {
		if part := strings.TrimSpace(arguments[i]); len(part) > 0 {
			if i == l && strings.HasSuffix(part, "...") {
				signature.IsVariadic = true
				part = part[:len(part)-3]
			}
			signature.Arguments = append(signature.Arguments, token.KindFromString(part))
		}
	}

	if size > bracketClose {
		signature.Returns = token.KindFromString(s[bracketClose+1:])
	}

	return &signature
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
	return a.lowerName == b.lowerName
}

func (sig *CommonSignature) MatchesArguments(args []token.Kind) bool {
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

func (sig *CommonSignature) sanitize() {
	sig.lowerName = strings.ToLower(sig.Name)
}

func (s *TaFunction) String() string {
	return s.CommonSignature.String()
}

func (a *TaFunction) Equal(b *TaFunction) bool {
	return a.CommonSignature.Equal(&b.CommonSignature)
}

func (s *TaFunction) MatchesArguments(args []token.Kind) bool {
	return s.CommonSignature.MatchesArguments(args)
}

func (s *TaTemplate) String() string {
	return s.CommonSignature.String()
}

func (a *TaTemplate) Equal(b *TaTemplate) bool {
	return a.CommonSignature.Equal(&b.CommonSignature)
}

func (s *TaTemplate) MatchesArguments(args []token.Kind) bool {
	return s.CommonSignature.MatchesArguments(args)
}
