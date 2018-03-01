//go:generate stringer -type=Kind

package block

import (
	"fmt"
	"strings"
	"time"

	"github.com/ericlagergren/decimal"
)

type Kind int

const (
	DecimalKind    Kind = 1 << iota
	StringKind     Kind = 1 << iota
	BoolKind       Kind = 1 << iota
	TimeKind       Kind = 1 << iota
	NullKind       Kind = 1 << iota
	ListKind       Kind = 1 << iota
	MapKind        Kind = 1 << iota
	BlockKind      Kind = 1 << iota
	AtomKind       Kind = DecimalKind | StringKind | BoolKind | TimeKind | NullKind
	CollectionKind Kind = BlockKind | ListKind | MapKind
	AnyKind        Kind = AtomKind | CollectionKind
)

type Block struct {
	String   string
	Decimal  *decimal.Big
	Bool     bool
	Time     time.Time
	Kind     Kind
	Children []*Block
	Keys     []string
}

func New(text string, children ...*Block) *Block {
	var b Block
	b.String = text
	if children == nil {
		b.Children = []*Block{}
	} else {
		b.Children = children
	}
	b.initValue(text)
	return &b
}

func NewDecimal(decimal *decimal.Big) *Block {
	var b Block
	b.Decimal = decimal
	b.Kind = DecimalKind
	b.String = b.Decimal.String()
	b.Children = []*Block{}
	return &b
}

func NewDecimalFromInt(i int64) *Block {
	var b Block
	b.Decimal = decimal.New(i, 0)
	b.Kind = DecimalKind
	b.String = b.Decimal.String()
	b.Children = []*Block{}
	return &b
}

func NewDecimalFromString(s string) *Block {
	var b Block
	var ok bool
	b.Children = []*Block{}
	b.Decimal = decimal.New(0, 0)
	b.Decimal, ok = b.Decimal.SetString(s)
	if !ok {
		b.Kind = NullKind
		return &b
	}
	b.Kind = DecimalKind
	b.String = b.Decimal.String()
	return &b
}

func NewBool(boolean bool) *Block {
	var b Block
	b.Bool = boolean
	b.Kind = BoolKind
	if boolean {
		b.String = "true"
	} else {
		b.String = "false"
	}
	b.Children = []*Block{}
	return &b
}

func NewTime(t time.Time) *Block {
	var b Block
	b.Time = t
	b.Kind = TimeKind
	b.String = b.Time.Format(time.RFC3339)
	b.Children = []*Block{}
	return &b
}

func NewString(str string) *Block {
	var b Block
	b.String = str
	b.Kind = StringKind
	b.Children = []*Block{}
	return &b
}

func NewNull() *Block {
	var b Block
	b.Kind = NullKind
	b.Children = []*Block{}
	return &b
}

func NewList(children ...*Block) *Block {
	var b Block
	if children == nil {
		b.Children = []*Block{}
	} else {
		b.Children = children
	}
	b.Kind = ListKind
	return &b
}

func NewMap(m map[string]*Block) *Block {
	var b Block
	b.Children = make([]*Block, len(m))
	b.Keys = make([]string, len(m))
	b.Kind = MapKind
	i := 0
	for k, v := range m {
		b.Keys[i] = k
		b.Children[i] = v
		i++
	}
	return &b
}

func (b *Block) IsEmpty() bool {
	return len(b.Children) == 0 && len(b.String) == 0
}

func (b *Block) IsDecimal() bool {
	return b.Kind == DecimalKind
}

func (b *Block) IsBool() bool {
	return b.Kind == BoolKind
}

func (b *Block) IsBlock() bool {
	return b.Kind == BlockKind
}

func (b *Block) IsTime() bool {
	return b.Kind == TimeKind
}

func (b *Block) IsString() bool {
	return b.Kind == StringKind
}

func (b *Block) IsNull() bool {
	return b.Kind == NullKind
}

func (b *Block) IsList() bool {
	return b.Kind == ListKind
}

func (b *Block) IsMap() bool {
	return b.Kind == MapKind
}

// Get an item from the map
func (b *Block) MapItem(key string) *Block {
	for i, k := range b.Keys {
		if key == k {
			return b.Children[i]
		}
	}
	return NewNull()
}

// Create the map
func (b *Block) Map() map[string]*Block {
	m := make(map[string]*Block)

	for i, key := range b.Keys {
		m[key] = b.Children[i]
	}
	return m
}

// todo: filterout invalid decimal types
func (b *Block) initValue(text string) {
	// only blocks could have children
	if len(b.Children) > 0 {
		b.Kind = BlockKind
		return
	}

	if len(b.String) > 0 {
		// is it a bool?
		if strings.EqualFold("true", text) {
			b.Bool = true
			b.Kind = BoolKind
			return
		} else if strings.EqualFold("false", text) {
			b.Bool = false
			b.Kind = BoolKind
			return
		}

		var err error
		b.Time, err = time.Parse(time.RFC3339, text)
		if err == nil {
			b.Kind = TimeKind
			return
		}

		var ok bool
		// try to parse it as a decimal
		b.Decimal, ok = decimal.New(0, 0).SetString(text)
		if ok {
			b.Kind = DecimalKind
			return
		}
		b.Kind = StringKind
	} else {
		b.Kind = BlockKind
	}
}

func (b *Block) Update(source *Block) {
	b.Kind = source.Kind
	switch b.Kind {
	case DecimalKind:
		b.Decimal = source.Decimal
	case BoolKind:
		b.Bool = source.Bool
	case TimeKind:
		b.Time = source.Time
	case MapKind:
		b.Keys = source.Keys
	}
	b.String = source.String
	b.Children = source.Children
}

func (b *Block) Stringify() string {
	text := b.String
	if l := len(b.Children); l > 0 {
		items := make([]string, l)
		for i, item := range b.Children {
			items[i] = item.Stringify()
		}
		text = fmt.Sprintf("%s %s", b.String, strings.Join(items, " "))
	}
	if b.IsBlock() {
		return fmt.Sprintf("(%s)", text)
	}
	return text
}

func Arguments(children []*Block) []Kind {
	types := make([]Kind, len(children))
	for i, child := range children {
		types[i] = child.Kind
	}
	return types
}

type BlockArguments []*Block

func (b BlockArguments) ToHumanReadable() string {
	arr := make([]string, len(b))
	for i, arg := range b {
		arr[i] = arg.Stringify()
	}
	return strings.Join(arr, ", ")
}
