package token

import (
	"fmt"
	"strings"
)

type Kind int

const (
	Decimal    Kind = 1 << iota
	String     Kind = 1 << iota
	Bool       Kind = 1 << iota
	Time       Kind = 1 << iota
	Null       Kind = 1 << iota
	List       Kind = 1 << iota
	Map        Kind = 1 << iota
	Token      Kind = 1 << iota
	Atom       Kind = Decimal | String | Bool | Time | Null
	Collection Kind = List | Map
	Any        Kind = Atom | Token | Collection
)

// cached for faster lookup
var kindStrings = map[Kind]string{
	Decimal:    strings.ToLower(Decimal.String()),
	String:     strings.ToLower(String.String()),
	Bool:       strings.ToLower(Bool.String()),
	Time:       strings.ToLower(Time.String()),
	Null:       strings.ToLower(Null.String()),
	List:       strings.ToLower(List.String()),
	Map:        strings.ToLower(Map.String()),
	Token:      strings.ToLower(Token.String()),
	Atom:       strings.ToLower(Atom.String()),
	Collection: strings.ToLower(Collection.String()),
	Any:        strings.ToLower(Any.String()),
}

func (k Kind) String() string {
	var kinds []string
	v := k
	for v != 0 {
		if v&Any == Any {
			kinds = append(kinds, "Any")
			v &= ^Any
		} else if v&Collection == Collection {
			kinds = append(kinds, "Collection")
			v &= ^Collection
		} else if v&Atom == Atom {
			kinds = append(kinds, "Atom")
			v &= ^Atom
		} else if v&Decimal == Decimal {
			kinds = append(kinds, "Decimal")
			v &= ^Decimal
		} else if v&String == String {
			kinds = append(kinds, "String")
			v &= ^String
		} else if v&Bool == Bool {
			kinds = append(kinds, "Bool")
			v &= ^Bool
		} else if v&Time == Time {
			kinds = append(kinds, "Time")
			v &= ^Time
		} else if v&Null == Null {
			kinds = append(kinds, "Null")
			v &= ^Null
		} else if v&List == List {
			kinds = append(kinds, "List")
			v &= ^List
		} else if v&Map == Map {
			kinds = append(kinds, "Map")
			v &= ^Map
		} else if v&Token == Token {
			kinds = append(kinds, "Token")
			v &= ^Token
		} else {
			kinds = append(kinds, fmt.Sprintf("Unknown(%d)", v))
			break
		}
	}
	return strings.Join(kinds, "|")
}

func KindFromString(s string) Kind {
	var v Kind
	parts := strings.Split(s, "|")
	for _, part := range parts {
		part = strings.ToLower(strings.TrimSpace(part))
		switch part {
		case kindStrings[Decimal]:
			v |= Decimal
		case kindStrings[String]:
			v |= String
		case kindStrings[Bool]:
			v |= Bool
		case kindStrings[Time]:
			v |= Time
		case kindStrings[Null]:
			v |= Null
		case kindStrings[List]:
			v |= List
		case kindStrings[Map]:
			v |= Map
		case kindStrings[Token]:
			v |= Token
		case kindStrings[Atom]:
			v |= Atom
		case kindStrings[Collection]:
			v |= Collection
		case kindStrings[Any]:
			v |= Any
		}
	}
	return v
}
