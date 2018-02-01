package talang

import (
	"strings"
	"unicode/utf8"

	"github.com/talon-one/talang/term"
)

// (= 1 2)

func Lex(str string) (term.Term, error) {
	// the first word is always the operation
	var children []term.Term
	var operation string

parse:
	strLen := len(str)
	start := 0
	i := 0

	for w := 0; i < strLen; i += w {
		r, width := utf8.DecodeRuneInString(str[i:])
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
					children = append(children, term.New(str[start:i]))
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
				return term.Term{}, err
			}
			if !nestedScope.IsEmpty() {
				children = append(children, nestedScope)
			}

			tokenString, rest := unquote(str[i:], r, r, '\\')
			if len(operation) == 0 {
				operation = tokenString
			} else {
				children = append(children, term.New(tokenString))
			}

			if rest == str {
				goto end
			}
			str = rest
			goto parse
		case 0x28: // bracket open
			tokenString, rest := unquote(str[i:], r, 0x29, utf8.RuneError)
			nestedScope, err := Lex(tokenString)
			if err != nil {
				return term.Term{}, err
			}
			if !nestedScope.IsEmpty() {
				children = append(children, nestedScope)
			}

			if rest == str {
				goto end
			}
			str = rest
			goto parse
		}

		w = width
	}
	if start < len(str) {
		if len(operation) == 0 {
			operation = str[start:]
		} else {
			children = append(children, term.New(str[start:]))
		}
	}

end:

	return term.New(operation, children...), nil
}

func unquote(str string, start, end, escape rune) (string, string) {
	// find which type of quote
	quoteRune, startRuneWidth := utf8.DecodeRuneInString(str)
	if quoteRune == utf8.RuneError || (quoteRune != start) {
		return "", str
	}
	_, endRuneWidth := utf8.DecodeRuneInString(string(end))

	l := len(str)
	if startRuneWidth >= l {
		return "", str
	}

	var nextQuotePos int
	var parts []string
	for i := startRuneWidth; i < len(str); {
		nextQuotePos = strings.IndexRune(str[i:], end)
		if nextQuotePos == -1 {
			return "", str
		}
		nextQuotePos += i
		if nextQuotePos >= l || nextQuotePos <= 0 {
			return "", ""
		}

		if escape != utf8.RuneError {
			previousRune, width := utf8.DecodeRuneInString(str[nextQuotePos-startRuneWidth:])
			if previousRune != escape {
				parts = append(parts, str[i:nextQuotePos])
				break
			}
			parts = append(parts, str[i:nextQuotePos-width]+string(end))
			nextQuotePos += width
		} else {
			// it might be nested
			s := i
			level := 0
			for s < nextQuotePos {
				pos := strings.IndexRune(str[s:nextQuotePos], start)
				if pos == -1 {
					break
				}
				s += pos + startRuneWidth
				level++
			}

			if level > 0 {
				for ; level >= 0; level-- {
					pos := strings.IndexRune(str[s:], end)
					if pos == -1 {
						break
					}
					s += pos + endRuneWidth
				}
				s -= endRuneWidth
				if s > 0 && s > i {
					parts = append(parts, str[i:s])
					nextQuotePos = s
				}
			} else {
				parts = append(parts, str[i:nextQuotePos])
				break
			}

		}

		i = nextQuotePos
	}

	unquoted := strings.Join(parts, "") //str[w:nextQuotePos]

	nextQuotePos++
	if nextQuotePos >= l {
		return unquoted, ""
	}
	return unquoted, str[nextQuotePos:]
}
