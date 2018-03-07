package unquote

import (
	"strings"
	"unicode/utf8"

	"github.com/talon-one/runes"
)

func Unquote(str string, quoteStart, quoteEnd string) (string, string) {
	if quoteStart == quoteEnd {
		return "", str
	}
	if !strings.HasPrefix(str, quoteStart) {
		return "", str
	}

	quoteStartLen := len(quoteStart)
	quoteEndLen := len(quoteEnd)

	pos := quoteStartLen

	endPos := strings.Index(str[pos:], quoteEnd)
	if endPos == -1 {
		return "", str
	}
	endPos += pos

	// no nesting
	firstNested := strings.Index(str[pos:], quoteStart)
	if firstNested == -1 {
		return str[quoteStartLen:endPos], str[endPos+quoteEndLen:]
	}

	// first nesting is after first quote
	firstNested += pos
	if firstNested > endPos {
		return str[quoteStartLen:endPos], str[endPos+quoteEndLen:]
	}

	// nesting is in our quote
	pos = firstNested
	for nestedStart := 0; nestedStart != -1; nestedStart = strings.Index(str[pos:], quoteStart) {
		pos += nestedStart
		unquoted, _ := Unquote(str[pos:], quoteStart, quoteEnd)
		consumed := len(unquoted)
		pos += quoteStartLen + consumed + quoteEndLen
	}

	endPosAfterNest := strings.Index(str[pos:], quoteEnd)
	if endPosAfterNest == -1 {
		return str[quoteStartLen:endPos], str[endPos+quoteEndLen:]
	}
	endPosAfterNest += pos

	return str[quoteStartLen:endPosAfterNest], str[endPosAfterNest+quoteEndLen:]
}

func EscapeUnquote(str string, quote string, escape string) (string, string) {
	unquoted, rest := escapeUnquoteRunes([]rune(str), []rune(quote), []rune(escape))
	return string(removeRune(unquoted, utf8.RuneError)), string(rest)
}

func removeRune(runes []rune, sep rune) []rune {
	if size := len(runes); size > 0 {
		v := make([]rune, 0, size)
		for i := 0; i < size; i++ {
			if r := runes[i]; r != sep {
				v = append(v, r)
			}
		}
		return v
	}
	return nil
}

func escapeUnquoteRunes(str []rune, quote []rune, escape []rune) ([]rune, []rune) {
	if !runes.HasPrefix(str, quote) {
		return nil, str
	}

	quoteLen := len(quote)
	escapeLen := len(escape)

	fullQuote := make([]rune, escapeLen+quoteLen)
	copy(fullQuote, escape)
	copy(fullQuote[escapeLen:], quote)

	pos := quoteLen

	endPos := runes.Index(str[pos:], quote)
	if endPos == -1 {
		return nil, str
	}
	endPos += pos

	// no nesting
	firstNested := runes.Index(str[pos:], fullQuote)
	if firstNested == -1 {
		return str[quoteLen:endPos], str[endPos+quoteLen:]
	}

	// first nesting is after first quote
	firstNested += pos
	if firstNested > endPos {
		return str[quoteLen:endPos], str[endPos+quoteLen:]
	}

	copyOfStr := runes.Copy(str)

	// nesting is in our quote
	pos = firstNested
	for nestedStart := 0; nestedStart != -1; nestedStart = runes.Index(str[pos:], fullQuote) {
		pos += nestedStart

		newQuote := make([]rune, escapeLen+quoteLen)
		copy(newQuote, escape)
		copy(newQuote[escapeLen:], quote)
		unquoted, _ := escapeUnquoteRunes(str[pos:], newQuote, escape)
		consumed := len(unquoted)
		if consumed > 0 {
			// mark it for removal
			str[pos] = utf8.RuneError
			str[pos+escapeLen+quoteLen+consumed] = utf8.RuneError
			pos += consumed + len(newQuote)*2
		} else {
			pos += escapeLen + quoteLen
		}
	}

	if size := len(str); pos >= size {
		return nil, copyOfStr
	}

	endPosAfterNest := runes.Index(str[pos:], quote)
	if endPosAfterNest == -1 {
		return str[quoteLen:endPos], str[endPos+quoteLen:]
	}
	endPosAfterNest += pos

	return str[quoteLen:endPosAfterNest], str[endPosAfterNest+quoteLen:]
}
