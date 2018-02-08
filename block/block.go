//go:generate stringer -type=Kind

package block

import (
	"fmt"
	"strings"

	"github.com/ericlagergren/decimal"
)

type Kind int

const (
	AnyKind     Kind = iota
	DecimalKind Kind = iota
	StringKind  Kind = iota
	BoolKind    Kind = iota
	BlockKind   Kind = iota
)

type Block struct {
	Text     string
	Decimal  *decimal.Big
	Bool     bool
	Kind     Kind
	Children []*Block
}

func (b *Block) IsEmpty() bool {
	return len(b.Children) == 0 && len(b.Text) == 0
}

func (b *Block) IsDecimal() bool {
	return b.Kind == DecimalKind
}

func (b *Block) IsBlock() bool {
	return b.Kind == BlockKind
}

func New(text string, children ...*Block) *Block {
	var b Block
	b.Text = text
	b.Children = children
	b.initValue(text)
	return &b
}

func NewDecimal(decimal *decimal.Big) *Block {
	var b Block
	b.Decimal = decimal
	b.Kind = DecimalKind
	b.Text = b.Decimal.String()
	return &b
}

func NewBool(boolean bool) *Block {
	var b Block
	b.Bool = boolean
	b.Kind = BoolKind
	if boolean {
		b.Text = "true"
	} else {
		b.Text = "false"
	}
	return &b
}
func NewString(str string) *Block {
	var b Block
	b.Text = str
	b.Kind = StringKind
	return &b
}

func (b *Block) initValue(text string) {
	// only blocks could have children
	if len(b.Children) > 0 {
		b.Kind = BlockKind
		return
	}

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

	var ok bool
	// try to parse it as a decimal
	b.Decimal, ok = decimal.New(0, 0).SetString(text)
	if ok == true {
		b.Kind = DecimalKind
		return
	}

	b.Kind = StringKind
}

func (b *Block) Update(source *Block) {
	b.Kind = source.Kind
	switch b.Kind {
	case DecimalKind:
		b.Decimal = source.Decimal
	case BoolKind:
		b.Bool = source.Bool
	}
	b.Text = source.Text
	b.Children = source.Children
}

func (b *Block) String() string {
	if l := len(b.Children); l > 0 {
		items := make([]string, l)
		for i, item := range b.Children {
			items[i] = item.String()
		}
		return fmt.Sprintf("(%s %s)", b.Text, strings.Join(items, " "))
	}
	return b.Text

}

func (b *Block) Arguments() []Kind {
	types := make([]Kind, len(b.Children))
	for i, child := range b.Children {
		types[i] = child.Kind
	}
	return types
}
