package term

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Term struct {
	Text      string
	Decimal   decimal.Decimal
	isDecimal bool
	Children  []Term
}

func (t *Term) IsEmpty() bool {
	return len(t.Children) == 0 && len(t.Text) == 0
}

func (t *Term) IsDecimal() bool {
	return t.isDecimal
}

func New(text string, children ...Term) (t Term) {
	t.Update(text, children...)
	return t
}

func (t *Term) Update(text string, children ...Term) {
	var err error
	t.Decimal, err = decimal.NewFromString(text)
	t.isDecimal = err == nil
	t.Text = text
	t.Children = children
}

func (t *Term) String() string {
	if l := len(t.Children); l > 0 {
		items := make([]string, l)
		for i, item := range t.Children {
			items[i] = item.String()
		}
		return fmt.Sprintf("(%s %s)", t.Text, strings.Join(items, " "))
	}
	return t.Text

}
