package lexer

import (
	"unicode/utf8"

	"github.com/talon-one/talang/block"
	"github.com/talon-one/talang/lexer/unquote"
)

func MustLex(str string) *block.TaToken {
	block, err := Lex(str)
	if err != nil {
		panic(err)
	}
	return block
}

func Lex(str string) (*block.TaToken, error) {
	// the first word is always the operation
	var children []*block.TaToken
	var operation string

parse:
	strLen := len(str)
	start := 0
	i := 0

	var r rune
	var width int
	for ; i < strLen; i += width {
		r, width = utf8.DecodeRuneInString(str[i:])
		switch r {
		case 0x09: // tab
			fallthrough
		case 0x0A: // LF
			fallthrough
		case 0x0D: // CR
			fallthrough
		case 0x20: // space
			if start < i {
				if len(operation) == 0 {
					operation = str[start:i]
				} else {
					children = append(children, block.New(str[start:i]))
				}
			}
			var j int
			for j = i + width; j < strLen; j += width {
				// decode next rune
				nextRune, _ := utf8.DecodeRuneInString(str[j:])
				if nextRune != r {
					break
				}
			}
			str = str[j:]
			goto parse
		case 0x22: // DoubleQuote
			fallthrough
		case 0x27: // SingleQoute
			nestedScope, err := Lex(str[start:i])
			if err != nil {
				return nil, err
			}
			if !nestedScope.IsEmpty() {
				children = append(children, nestedScope)
			}

			tokenString, rest := unquote.EscapeUnquote(str[i:], string(r), `\`)
			if len(operation) == 0 {
				operation = tokenString
			} else {
				children = append(children, block.NewString(tokenString))
			}

			if rest == str {
				goto end
			}
			str = rest
			goto parse
		case 0x28: // bracket open
			tokenString, rest := unquote.Unquote(str[i:], "(", ")")
			nestedScope, err := Lex(tokenString)
			if err != nil {
				return nil, err
			}
			if !nestedScope.IsEmpty() {
				nestedScope.Kind = block.Token
				children = append(children, nestedScope)
			}

			if rest == str {
				goto end
			}
			str = rest
			goto parse
		}
	}
	if start < len(str) {
		if len(operation) == 0 {
			operation = str[start:]
		} else {
			children = append(children, block.New(str[start:]))
		}
	}

end:

	return block.New(operation, children...), nil
}
