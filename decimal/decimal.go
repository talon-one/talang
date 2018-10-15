package decimal

// this is a wrapper class for decimal
// so it is possible to change the internal decimal library later with another one

import (
	"fmt"

	"github.com/ericlagergren/decimal"
)

type Decimal struct {
	dec *decimal.Big
}

func newFromDec(d *Decimal) *Decimal {
	return &Decimal{
		dec: d.dec.Copy(d.dec),
	}
}

func NewFromInt(i int64) *Decimal {
	return &Decimal{
		dec: decimal.New(i, 0),
	}
}

func NewFromString(s string) *Decimal {
	var d Decimal
	var ok bool
	d.dec = decimal.New(0, 0)
	d.dec, ok = d.dec.SetString(s)
	if !ok {
		return nil
	}
	return &d
}

// NewFromFloat creates a new decimal block from a float (REMEMBER float64 is not exact! use with care)
func NewFromFloat(i float64) *Decimal {
	var d Decimal
	d.dec = decimal.New(0, 0)
	d.dec = d.dec.SetFloat64(i)
	return &d
}
func (d *Decimal) Equal(a *Decimal) bool {
	if a == nil || a.dec == nil {
		return false
	}
	return d.Compare(a) == 0
}

func (d *Decimal) Compare(a *Decimal) int {
	if a == nil || a.dec == nil {
		return -1
	}
	return d.dec.Cmp(a.dec)
}

func (d *Decimal) String() string {
	if d.dec == nil {
		return ""
	}
	return d.dec.String()
}

func (d *Decimal) Int8() (int8, error) {
	i, ok := d.dec.Int64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an int8", d.String())
	}
	return int8(i), nil
}

func (d *Decimal) Int16() (int16, error) {
	i, ok := d.dec.Int64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an int16", d.String())
	}
	return int16(i), nil
}

func (d *Decimal) Int32() (int32, error) {
	i, ok := d.dec.Int64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an int32", d.String())
	}
	return int32(i), nil
}

func (d *Decimal) Int64() (int64, error) {
	i, ok := d.dec.Int64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an int64", d.String())
	}
	return i, nil
}

func (d *Decimal) Uint8() (uint8, error) {
	i, ok := d.dec.Uint64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an uint8", d.String())
	}
	return uint8(i), nil
}

func (d *Decimal) Uint16() (uint16, error) {
	i, ok := d.dec.Uint64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an Uint16", d.String())
	}
	return uint16(i), nil
}

func (d *Decimal) Uint32() (uint32, error) {
	i, ok := d.dec.Uint64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an uint32", d.String())
	}
	return uint32(i), nil
}

func (d *Decimal) Uint64() (uint64, error) {
	i, ok := d.dec.Uint64()
	if !ok {
		return 0, fmt.Errorf("`%s' not anu int64", d.String())
	}
	return i, nil
}

func (d *Decimal) Int() (int, error) {
	i, ok := d.dec.Int64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an int64", d.String())
	}
	return int(i), nil
}

func (d *Decimal) Uint() (uint, error) {
	i, ok := d.dec.Uint64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an int64", d.String())
	}
	return uint(i), nil
}

func (d *Decimal) Float32() (float32, error) {
	i, ok := d.dec.Float64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an float32", d.String())
	}
	return float32(i), nil
}

func (d *Decimal) Float64() (float64, error) {
	i, ok := d.dec.Float64()
	if !ok {
		return 0, fmt.Errorf("`%s' not an float32", d.String())
	}
	return i, nil
}
