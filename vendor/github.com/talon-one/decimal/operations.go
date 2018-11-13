package decimal

import (
	"strings"

	"github.com/ericlagergren/decimal/math"
)

// Cmp compares n to the decimal instance
func (dec Decimal) Cmp(n Decimal) int {
	return dec.native.Cmp(n.native)
}

// Cmp compares a to b
func Cmp(a Decimal, b Decimal) int {
	return a.Cmp(b)
}

// Equals returns true if n has the same value as the decimal instance
func (dec Decimal) Equals(n Decimal) bool {
	return dec.Cmp(n) == 0
}

// Equals returns true if a has the same value as b
func Equals(a Decimal, b Decimal) bool {
	return a.Equals(b)
}

// EqualsInterface returns true if v is an decimal and has the same value as the decimal instance
func (dec Decimal) EqualsInterface(v interface{}) bool {
	switch x := v.(type) {
	case Decimal:
		return dec.Cmp(x) == 0
	case *Decimal:
		return dec.Cmp(*x) == 0
	}
	return false
}

// EqualsInterface returns true if v is an decimal and has the same value as d
func EqualsInterface(d Decimal, v interface{}) bool {
	return d.EqualsInterface(v)
}

// Add adds n to the decimal instance
func (dec Decimal) Add(n Decimal) Decimal {
	dec.native.Add(dec.native, n.native)
	return Decimal{dec.native}
}

// Add adds a to b and returns a new decimal instance
// a and b will not be modified
func Add(a Decimal, b Decimal) Decimal {
	d := NewFromDecimal(a)
	return d.Add(b)
}

// Sub substracts n to the decimal instance
func (dec Decimal) Sub(n Decimal) Decimal {
	dec.native.Sub(dec.native, n.native)
	return Decimal{dec.native}
}

// Sub substracts b from a and returns a new decimal instance
// a and b will not be modified
func Sub(a Decimal, b Decimal) Decimal {
	d := NewFromDecimal(a)
	return d.Sub(b)
}

// Div divides n on the decimal instance
func (dec Decimal) Div(n Decimal) Decimal {
	dec.native.Quo(dec.native, n.native)
	return Decimal{dec.native}
}

// Div divides b from a and returns a new decimal instance
// a and b will not be modified
func Div(a Decimal, b Decimal) Decimal {
	d := NewFromDecimal(a)
	return d.Div(b)
}

// Mul multiplies n to the decimal instance
func (dec Decimal) Mul(n Decimal) Decimal {
	dec.native.Mul(dec.native, n.native)
	return Decimal{dec.native}
}

// Mul multiplies a to b and returns a new decimal instance
// a and b will not be modified
func Mul(a Decimal, b Decimal) Decimal {
	d := NewFromDecimal(a)
	return d.Mul(b)
}

// Mod modulos n on the decimal instance
func (dec Decimal) Mod(n Decimal) Decimal {
	dec.native.Rem(dec.native, n.native)
	return Decimal{dec.native}
}

// Mod modulos b on a and returns a new decimal instance
// a and b will not be modified
func Mod(a Decimal, b Decimal) Decimal {
	d := NewFromDecimal(a)
	return d.Mod(b)
}

// Floor rounds the instance down to the next whole number
func (dec Decimal) Floor() Decimal {
	math.Floor(dec.native, dec.native)
	return Decimal{dec.native}
}

// Floor rounds d down to the next whole number and returns it as a new instance
// d will not be modified
func Floor(a Decimal) Decimal {
	d := NewFromDecimal(a)
	return d.Floor()
}

// Ceil rounds the instance up to the next whole number
func (dec Decimal) Ceil() Decimal {
	math.Ceil(dec.native, dec.native)
	return Decimal{dec.native}
}

// Ceil rounds d up to the next whole number and returns it as a new instance
// d will not be modified
func Ceil(a Decimal) Decimal {
	d := NewFromDecimal(a)
	return d.Ceil()
}

// Round rounds the instance to the specific digits
func (dec Decimal) Round(digits int) Decimal {
	dec.native.Round(digits)
	return Decimal{dec.native}
}

// Round rounds d to the specific digits and returns it as a new instance
// d will not be modified
func Round(a Decimal, digits int) Decimal {
	d := NewFromDecimal(a)
	return d.Round(digits)
}

// Truncate truncates the instance to the specific digits
func (dec Decimal) Truncate(digits int) Decimal {
	parts := strings.SplitN(dec.native.String(), ".", 2)
	if len(parts) <= 1 {
		v, _ := NewFromString(parts[0])
		dec.native.Copy(v.native)
		return v
	}
	if digits > len(parts[1])-1 {
		digits = len(parts[1])
	}
	v, _ := NewFromString(parts[0] + "." + parts[1][:digits])
	dec.native.Copy(v.native)
	return v
}

// Truncate truncates d to the specific digits and returns it as a new instance
// d will not be modified
func Truncate(a Decimal, digits int) Decimal {
	d := NewFromDecimal(a)
	return d.Truncate(digits)
}
