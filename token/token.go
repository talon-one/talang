package token

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/ericlagergren/decimal"
)

type TaToken struct {
	String   string
	Decimal  *decimal.Big
	Bool     bool
	Time     time.Time
	Kind     Kind
	Children []*TaToken
	Keys     []string
}

func New(text string, children ...*TaToken) *TaToken {
	var b TaToken
	b.String = text
	if children == nil {
		b.Children = []*TaToken{}
	} else {
		b.Children = children
	}
	b.initValue(text)
	return &b
}

func NewDecimal(decimal *decimal.Big) *TaToken {
	var b TaToken
	b.Decimal = decimal
	b.Kind = Decimal
	b.String = b.Decimal.String()
	b.Children = []*TaToken{}
	return &b
}

func NewDecimalFromInt(i int64) *TaToken {
	var b TaToken
	b.Decimal = decimal.New(i, 0)
	b.Kind = Decimal
	b.String = b.Decimal.String()
	b.Children = []*TaToken{}
	return &b
}

func NewDecimalFromString(s string) *TaToken {
	var b TaToken
	var ok bool
	b.Children = []*TaToken{}
	b.Decimal = decimal.New(0, 0)
	b.Decimal, ok = b.Decimal.SetString(s)
	if !ok {
		b.Kind = Null
		return &b
	}
	b.Kind = Decimal
	b.String = b.Decimal.String()
	return &b
}

// NewDecimalFromFloat creates a new decimal block from a float (REMEMBER float64 is not exact! use with care)
func NewDecimalFromFloat(i float64) *TaToken {
	var b TaToken
	b.Decimal = decimal.New(0, 0)
	b.Decimal = b.Decimal.SetFloat64(i)
	b.Kind = Decimal
	b.String = b.Decimal.String()
	b.Children = []*TaToken{}
	return &b
}

func NewBool(boolean bool) *TaToken {
	var b TaToken
	b.Bool = boolean
	b.Kind = Boolean
	if boolean {
		b.String = "true"
	} else {
		b.String = "false"
	}
	b.Children = []*TaToken{}
	return &b
}

func NewTime(t time.Time) *TaToken {
	var b TaToken
	b.Time = t
	b.Kind = Time
	b.String = b.Time.Format(time.RFC3339)
	b.Children = []*TaToken{}
	return &b
}

func NewString(str string) *TaToken {
	var b TaToken
	b.String = str
	b.Kind = String
	b.Children = []*TaToken{}
	return &b
}

func NewNull() *TaToken {
	var b TaToken
	b.Kind = Null
	b.Children = []*TaToken{}
	return &b
}

func NewList(children ...*TaToken) *TaToken {
	var b TaToken
	if children == nil {
		b.Children = []*TaToken{}
	} else {
		b.Children = children
	}
	b.Kind = List
	return &b
}

func NewMap(m map[string]*TaToken) *TaToken {
	var b TaToken
	b.Children = make([]*TaToken, len(m))
	b.Keys = make([]string, len(m))
	b.Kind = Map
	i := 0
	for k, v := range m {
		b.Keys[i] = k
		b.Children[i] = v
		i++
	}
	return &b
}

func NewToken(text string, children ...*TaToken) *TaToken {
	var b TaToken
	b.String = text
	b.Kind = Token
	if children == nil {
		b.Children = []*TaToken{}
	} else {
		b.Children = children
	}
	return &b
}

func (b *TaToken) IsEmpty() bool {
	return len(b.Children) == 0 && len(b.String) == 0
}

func (b *TaToken) IsDecimal() bool {
	return b.Kind == Decimal
}

func (b *TaToken) IsBool() bool {
	return b.Kind == Boolean
}

func (b *TaToken) IsBlock() bool {
	return b.Kind == Token
}

func (b *TaToken) IsTime() bool {
	return b.Kind == Time
}

func (b *TaToken) IsString() bool {
	return b.Kind == String
}

func (b *TaToken) IsNull() bool {
	return b.Kind == Null
}

func (b *TaToken) IsList() bool {
	return b.Kind == List
}

func (b *TaToken) IsMap() bool {
	return b.Kind == Map
}

// Get an item from the map
func (b *TaToken) MapItem(key string) *TaToken {
	for i, k := range b.Keys {
		if key == k {
			return b.Children[i]
		}
	}
	return NewNull()
}

// Set an item in the map
func (b *TaToken) SetMapItem(key string, value *TaToken) {
	for i, k := range b.Keys {
		if key == k {
			b.Children[i] = value
			return
		}
	}
	b.Keys = append(b.Keys, key)
	b.Children = append(b.Children, value)
}

// Create the map
func (b *TaToken) Map() map[string]*TaToken {
	m := make(map[string]*TaToken)

	for i, key := range b.Keys {
		m[key] = b.Children[i]
	}
	return m
}

func (b *TaToken) initValue(text string) {
	// only blocks could have children
	if len(b.Children) > 0 {
		b.Kind = Token
		return
	}

	if len(b.String) > 0 {
		// is it a bool?
		if strings.EqualFold("true", text) {
			b.Bool = true
			b.Kind = Boolean
			return
		} else if strings.EqualFold("false", text) {
			b.Bool = false
			b.Kind = Boolean
			return
		}

		// is it a decimal?
		if isDecimal(text) {
			var ok bool
			// try to parse it as a decimal
			b.Decimal, ok = decimal.New(0, 0).SetString(text)
			if ok {
				b.Kind = Decimal
				return
			}
		}

		// is it a time?
		var err error
		b.Time, err = time.Parse(time.RFC3339, text)
		if err == nil {
			b.Kind = Time
			return
		}

		b.Kind = String
	} else {
		b.Kind = Token
	}
}

func isDecimal(s string) bool {
	if len(s) <= 0 {
		return false
	}
	runes := []rune(s)

	i := 0
	if runes[0] == '+' || runes[0] == '-' {
		i++
	}

	gotDot := false
	for ; i < len(runes); i++ {
		if runes[i] == '.' {
			if gotDot {
				return false
			}
			gotDot = true
			continue
		}
		if !unicode.IsNumber(runes[i]) {
			return false
		}
	}

	return true
}

// Copy creates a copy of the block
func Copy(dst *TaToken, src *TaToken) {
	if src == nil {
		return
	}
	if dst == nil {
		return
	}
	dst.Kind = src.Kind
	switch dst.Kind {
	case Decimal:
		dst.Decimal = src.Decimal
	case Boolean:
		dst.Bool = src.Bool
	case Time:
		dst.Time = src.Time
	case Map:
		dst.Keys = make([]string, len(src.Keys))
		copy(dst.Keys, src.Keys)
	}
	dst.String = src.String
	dst.Children = make([]*TaToken, len(src.Children))
	for i, child := range src.Children {
		dst.Children[i] = new(TaToken)
		Copy(dst.Children[i], child)
	}
}

func (b *TaToken) Stringify() string {
	var builder strings.Builder
	var children []string
	if l := len(b.Children); l > 0 {
		children = make([]string, l)
		for i, item := range b.Children {
			children[i] = item.Stringify()
		}
	}
	if b.IsBlock() {
		builder.WriteString("(")
		if len(b.String) > 0 {
			builder.WriteString(b.String)
			if len(children) > 0 {
				builder.WriteString(" ")
			}
		}
		builder.WriteString(strings.Join(children, " "))
		builder.WriteString(")")
	} else if b.IsList() {
		builder.WriteString("[")
		builder.WriteString(strings.Join(children, ", "))
		builder.WriteString("]")
	} else if b.IsMap() {
		builder.WriteString("{")
		var keys []string
		if l := len(b.Keys); l > 0 {
			keys = make([]string, l)
			for i, item := range b.Keys {
				keys[i] = fmt.Sprintf("%s:%s", item, children[i])
			}
		}
		builder.WriteString(strings.Join(keys, ", "))
		builder.WriteString("}")
	} else if len(b.String) > 0 {
		if b.IsString() {
			builder.WriteRune('"')
		}
		builder.WriteString(b.String)
		if b.IsString() {
			builder.WriteRune('"')
		}
	}

	return builder.String()
}

func (b *TaToken) Equal(a *TaToken) bool {
	if a == nil || a.Kind != b.Kind {
		return false
	}

	switch a.Kind {
	case Decimal:
		return a.Decimal.Cmp(b.Decimal) == 0
	case String:
		return a.String == b.String
	case Boolean:
		return a.Bool == b.Bool
	case Time:
		return a.Time.Equal(b.Time)
	case Null:
		return true
	case Map:
		if len(a.Keys) != len(b.Keys) {
			return false
		}

		for _, key := range a.Keys {
			if !a.MapItem(key).Equal(b.MapItem(key)) {
				return false
			}
		}
		return true
	case List:
		fallthrough
	case Token:
		if len(a.Children) != len(b.Children) {
			return false
		}
		for i, item := range a.Children {
			if !item.Equal(b.Children[i]) {
				return false
			}
		}
		return true
	}
	return false
}

func Arguments(children []*TaToken) []Kind {
	types := make([]Kind, len(children))
	for i, child := range children {
		types[i] = child.Kind
	}
	return types
}

type TokenArguments []*TaToken

func (b TokenArguments) ToHumanReadable() string {
	arr := make([]string, len(b))
	for i, arg := range b {
		arr[i] = arg.Stringify()
	}
	return strings.Join(arr, ", ")
}
func (b TokenArguments) Len() int           { return len(b) }
func (b TokenArguments) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b TokenArguments) Less(i, j int) bool { return strings.Compare(b[i].String, b[j].String) < 0 }

func (k *Kind) UnmarshalJSON(b []byte) (err error) {
	*k = KindFromString(strings.Trim(string(b), `"`))
	return nil
}

func (k *Kind) MarshalJSON() ([]byte, error) {
	var builder strings.Builder
	builder.WriteString(`"`)
	builder.WriteString(k.String())
	builder.WriteString(`"`)
	return []byte(builder.String()), nil
}
