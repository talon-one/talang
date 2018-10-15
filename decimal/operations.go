package decimal

import (
	"github.com/ericlagergren/decimal"
)

func Add(a, b *Decimal) *Decimal {
	r := newFromDec(a)
	r.Add(b)
	return r
}

func Sub(a, b *Decimal) *Decimal {
	r := newFromDec(a)
	r.Sub(b)
	return r
}

func Mul(a, b *Decimal) *Decimal {
	r := newFromDec(a)
	r.Mul(b)
	return r
}

func Div(a, b *Decimal) *Decimal {
	r := newFromDec(a)
	r.Div(b)
	return r
}

func Mod(a, b *Decimal) *Decimal {
	r := newFromDec(a)
	r.Mod(b)
	return r
}

func Floor(a *Decimal) *Decimal {
	r := newFromDec(a)
	r.Floor()
	return r
}

func Ceil(a *Decimal) *Decimal {
	r := newFromDec(a)
	r.Ceil()
	return r
}

func (d *Decimal) Add(a *Decimal) {
	d.dec = d.dec.Add(d.dec, a.dec)
}

func (d *Decimal) Sub(a *Decimal) {
	d.dec = d.dec.Sub(d.dec, a.dec)
}

func (d *Decimal) Mul(a *Decimal) {
	d.dec = d.dec.Mul(d.dec, a.dec)
}

func (d *Decimal) Div(a *Decimal) {
	d.dec = d.dec.Quo(d.dec, a.dec)
}

func (d *Decimal) Mod(a *Decimal) {
	d.dec = d.dec.Rem(d.dec, a.dec)
}

func (d *Decimal) Floor() {
	ctx := decimal.Context{Precision: d.dec.Context.Precision}
	if d.dec.Signbit() {
		ctx.RoundingMode = decimal.AwayFromZero
	} else {
		ctx.RoundingMode = decimal.ToZero
	}
	d.dec = ctx.RoundToInt(d.dec)
}

func (d *Decimal) Ceil() {
	ctx := decimal.Context{Precision: d.dec.Context.Precision}
	if d.dec.Signbit() {
		ctx.RoundingMode = decimal.ToZero
	} else {
		ctx.RoundingMode = decimal.AwayFromZero
	}
	d.dec = ctx.RoundToInt(d.dec)
}
