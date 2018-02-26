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
	DecimalKind Kind = 1 << iota
	StringKind  Kind = 1 << iota
	BoolKind    Kind = 1 << iota
	TimeKind    Kind = 1 << iota
	BlockKind   Kind = 1 << iota
	NullKind    Kind = 1 << iota
	AtomKind    Kind = DecimalKind | StringKind | BoolKind | TimeKind | NullKind
	AnyKind     Kind = AtomKind | BlockKind
)

type Block struct {
	String   string
	Decimal  *decimal.Big
	Bool     bool
	Time     time.Time
	Kind     Kind
	Children []*Block
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
