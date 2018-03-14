package decimal

import (
	"bytes"
	"encoding"
	"fmt"
	"io"
	"math"
	"math/big"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/ericlagergren/decimal/internal/arith"
	"github.com/ericlagergren/decimal/internal/arith/checked"
	"github.com/ericlagergren/decimal/internal/c"
)

// Big is a floating-point, arbitrary-precision decimal.
//
// It is represented as a number and a scale. If the scale is >= 0, it indicates
// the number of decimal digits after the radix. Otherwise, the number is
// multiplied by 10 to the power of the negation of the scale. More formally,
//
//   Big = number × 10**-scale
//
// with MinScale <= scale <= MaxScale. A Big may also be ±0, ±Infinity, or ±NaN
// (either quiet or signaling). Non-NaN Big values are ordered, defined as the
// result of x.Cmp(y).
//
// Additionally, each Big value has a contextual object which governs arithmetic
// operations.
type Big struct {
	// Context is the decimal's unique contextual object.
	Context Context

	// unscaled is only used if the decimal is too large to fit in compact.
	unscaled big.Int

	// compact is use if the value fits into an uint64. The scale does not
	// affect whether this field is used; typically, if a decimal has <= 20
	// digits this field will be used.
	compact uint64

	// exp is the negated scale, meaning
	//
	//   number × 10**exp = number ×  10**-scale
	exp int

	// precision is the current precision.
	precision int

	// form indicates whether a decimal is a finite number, an infinity, or a
	// NaN value and whether it's signed or not.
	form form
}

// form indicates whether a decimal is a finite number, an infinity, or a nan
// value and whether it's signed or not.
type form uint8

const (
	// Particular bits:
	//
	// 0: sign bit
	// 1: infinity
	// 2: signaling nan
	// 3: quiet nan
	// 4-7: unused

	finite form = 0 // default, all zeros; do not re-order this constant.

	signbit form = 1 << 0 // do not assign this; used to check for signedness.

	pinf form = 1 << 1         // may compare with ==, &, etc.
	ninf form = pinf | signbit // may compare with ==, &, etc.
	inf  form = pinf           // do not assign this; used to check for either infinity.

	snan  form = 1 << 2         // compare with bitwise & only due to ssnan
	qnan  form = 1 << 3         // compare with bitwise & only due to sqnan
	ssnan form = snan | signbit // primarily for printing, signbit
	sqnan form = qnan | signbit // primarily for printing, signbit
	nan   form = snan | qnan    // do not assign this; used to check for either NaN.

	special = inf | nan // do not assign this; used to check for a special value.
)

func (f form) String() string {
	// GDA versions. Go needs to be handled manually.
	switch f {
	case finite:
		return "finite"
	case finite | signbit:
		return "-finite"
	case snan:
		return "sNaN"
	case snan | signbit:
		return "-sNaN"
	case qnan:
		return "NaN"
	case qnan | signbit:
		return "-NaN"
	case pinf:
		return "Infinity"
	case ninf:
		return "-Infinity"
	default:
		return fmt.Sprintf("unknown form: %0.8b", f)
	}
}

// Payload is a NaN value's payload.
type Payload uint64

const (
	addinfinf Payload = iota + 1
	mul0inf
	quo00
	quoinfinf
	quantinf
	quantminmax
	quantprec
	subinfinf
	absvalue
	addition
	comparison
	multiplication
	negation
	division
	quantization
	subtraction
	quorem_
	reminfy
	remx0
	quotermexp
	invctxpltz
	invctxpgtu
	invctxrmode
	invctxomode
	invctxsltu
	invctxsgtu
	reduction
	quointprec
	remprec
)

var payloads = [...]string{
	addinfinf:      "addition of infinities with opposing signs",
	mul0inf:        "multiplication of zero with infinity",
	quo00:          "division of zero by zero",
	quoinfinf:      "division of infinity by infinity",
	quantinf:       "quantization of an infinity",
	quantminmax:    "quantization exceeds minimum or maximum scale",
	quantprec:      "quantization exceeds working precision",
	subinfinf:      "subtraction of infinities with opposing signs",
	absvalue:       "absolute value of NaN",
	addition:       "addition with NaN as an operand",
	comparison:     "comparison with NaN as an operand",
	multiplication: "multiplication with NaN as an operand",
	negation:       "negation with NaN as an operand",
	division:       "division with NaN as an operand",
	quantization:   "quantization with NaN as an operand",
	subtraction:    "subtraction with NaN as an operand",
	quorem_:        "integer division or remainder has too many digits",
	reminfy:        "remainder of infinity",
	remx0:          "remainder by zero",
	quotermexp:     "division with unlimited precision has a non-terminating decimal expansion",
	invctxpltz:     "operation with a precision less than zero",
	invctxpgtu:     "operation with a precision greater than MaxPrecision",
	invctxrmode:    "operation with an invalid RoundingMode",
	invctxomode:    "operation with an invalid OperatingMode",
	invctxsltu:     "operation with a scale lesser than MinScale",
	invctxsgtu:     "operation with a scale greater than MaxScale",
	reduction:      "reduction with NaN as an operand",
	quointprec:     "result of integer division was larger than the desired precision",
	remprec:        "result of remainder operation was larger than the desired precision",
}

func (p Payload) String() string {
	if p < Payload(len(payloads)) {
		return payloads[p]
	}
	return ""
}

// An ErrNaN is used when a decimal operation would lead to a NaN under IEEE-754
// rules. An ErrNaN implements the error interface.
type ErrNaN struct{ Msg string }

func (e ErrNaN) Error() string { return e.Msg }

var _ error = ErrNaN{}

// CheckNaNs checks if either x or y is NaN. If so, it follows the rules of NaN
// handling set forth in the GDA specification. The second argument, y, may be
// nil. It returns true if either condition is a NaN.
func (z *Big) CheckNaNs(x, y *Big) bool {
	return z.invalidContext(z.Context) || z.checkNaNs(x, y, 0)
}

func (z *Big) checkNaNs(x, y *Big, op Payload) bool {
	var yform form
	if y != nil {
		yform = y.form
	}
	f := (x.form | yform) & nan
	if f == 0 {
		return false
	}

	form := qnan
	var cond Condition
	if f&snan != 0 {
		cond = InvalidOperation
		if x.form&snan != 0 {
			form |= (x.form & signbit)
		} else {
			form |= (y.form & signbit)
		}
	} else if x.form&nan != 0 {
		form |= (x.form & signbit)
	} else {
		form |= (y.form & signbit)
	}
	z.setNaN(cond, form, op)
	return true
}

func (z *Big) xflow(exp int, over, neg bool) *Big {
	// over == overflow
	// neg == intermediate result < 0
	if over {
		// TODO(eric): actually choose the largest finite number in the current
		// precision. This is legacy now.
		//
		// NOTE(eric): in some situations, the decimal library tells us to set
		// z to "the largest finite number that can be represented in the
		// current precision..." Use signed Infinity instead.
		//
		// Because of the logic above, every rounding mode works out to the
		// following.
		if neg {
			z.form = ninf
		} else {
			z.form = pinf
		}
		z.Context.Conditions |= Overflow | Inexact | Rounded
		return z
	}

	var sign form
	if neg {
		sign = signbit
	}
	z.setZero(sign, exp)
	z.Context.Conditions |= Underflow | Inexact | Rounded | Subnormal
	return z
}

// These methods are here to prevent typos.

func (x *Big) isCompact() bool  { return x.compact != c.Inflated }
func (x *Big) isInflated() bool { return !x.isCompact() }
func (x *Big) isSpecial() bool  { return x.form&(inf|nan) != 0 }

func (x *Big) adjusted() int { return (x.exp + x.Precision()) - 1 }
func (c Context) etiny() int { return MinScale - (precision(c) - 1) }

// Abs sets z to the absolute value of x and returns z.
func (z *Big) Abs(x *Big) *Big {
	if debug {
		x.validate()
	}
	if !z.invalidContext(z.Context) && !z.checkNaNs(x, x, absvalue) {
		z.Context.round(z.copyAbs(x))
	}
	return z
}

// Add sets z to x + y and returns z.
func (z *Big) Add(x, y *Big) *Big { return z.Context.Add(z, x, y) }

// Class returns the ``class'' of x, which is one of the following:
//
//  sNaN
//  NaN
//  -Infinity
//  -Normal
//  -Subnormal
//  -Zero
//  +Zero
//  +Subnormal
//  +Normal
//  +Infinity
//
func (x *Big) Class() string {
	if x.IsNaN(0) {
		if x.IsNaN(+1) {
			return "NaN"
		}
		return "sNaN"
	}
	if x.Signbit() {
		if x.IsInf(0) {
			return "-Infinity"
		}
		if x.compact == 0 {
			return "-Zero"
		}
		if x.IsSubnormal() {
			return "-Subnormal"
		}
		return "-Normal"
	}
	if x.IsInf(0) {
		return "+Infinity"
	}
	if x.compact == 0 {
		return "+Zero"
	}
	if x.IsSubnormal() {
		return "+Subnormal"
	}
	return "+Normal"
}

// Cmp compares x and y and returns:
//
//   -1 if x <  y
//    0 if x == y
//   +1 if x >  y
//
// It does not modify x or y. The result is undefined if either x or y are NaN.
// For an abstract comparison with NaN values, see misc.CmpTotal.
func (x *Big) Cmp(y *Big) int { return cmp(x, y, false) }

// CmpAbs compares |x| and |y| and returns:
//
//   -1 if |x| <  |y|
//    0 if |x| == |y|
//   +1 if |x| >  |y|
//
// It does not modify x or y. The result is undefined if either x or y are NaN.
// For an abstract comparison with NaN values, see misc.CmpTotalAbs.
func (x *Big) CmpAbs(y *Big) int { return cmp(x, y, true) }

// cmp is the implementation for both Cmp and CmpAbs.
func cmp(x, y *Big, abs bool) int {
	if debug {
		x.validate()
		y.validate()
	}

	if x == y {
		return 0
	}

	// NaN cmp x
	// z cmp NaN
	// NaN cmp NaN
	if (x.form|y.form)&nan != 0 {
		return 0
	}

	// Fast path: Catches non-finite forms like zero and ±Inf, possibly signed.
	xs := x.ord(abs)
	ys := y.ord(abs)
	if xs != ys {
		if xs > ys {
			return +1
		}
		return -1
	}
	switch xs {
	case 0, +2, -2:
		return 0
	default:
		r := cmpabs(x, y)
		if xs < 0 && !abs {
			r = -r
		}
		return r
	}
}

func cmpabs(x, y *Big) int {
	// Same scales means we can compare straight across.
	if x.exp == y.exp {
		if x.isCompact() {
			if y.isCompact() {
				return arith.Cmp(x.compact, y.compact)
			}
			return -1 // y.isInflateed
		}
		if y.isCompact() {
			return +1 // !x.isCompact
		}
		return x.unscaled.CmpAbs(&y.unscaled)
	}

	// Signs are the same and the scales differ. Compare the lengths of their
	// integral parts; if they differ in length one number is larger.
	// E.g., 1234.01
	//        123.011
	xl := x.adjusted()
	yl := y.adjusted()

	if xl != yl {
		if xl < yl {
			return -1
		}
		return +1
	}

	diff := int64(x.exp) - int64(y.exp)
	shift := uint64(arith.Abs(diff))
	if arith.Safe(shift) && x.isCompact() && y.isCompact() {
		p, _ := arith.Pow10(shift)
		if diff < 0 {
			return arith.AbsCmp128(x.compact, y.compact, p)
		}
		return -arith.AbsCmp128(y.compact, x.compact, p)
	}

	xw, yw := x.unscaled.Bits(), y.unscaled.Bits()
	if x.isCompact() {
		xw = arith.Words(x.compact)
	}
	if y.isCompact() {
		yw = arith.Words(y.compact)
	}

	var tmp big.Int
	if diff < 0 {
		yw = checked.MulBigPow10(&tmp, tmp.SetBits(copybits(yw)), shift).Bits()
	} else {
		xw = checked.MulBigPow10(&tmp, tmp.SetBits(copybits(xw)), shift).Bits()
	}
	return arith.CmpBits(xw, yw)
}

// Copy sets z to a copy of x and returns z.
func (z *Big) Copy(x *Big) *Big {
	if debug {
		x.validate()
	}
	if z != x {
		sign := x.form & signbit
		z.copyAbs(x)
		z.form |= sign
	}
	return z
}

// copyAbs sets z to a copy of |x| and returns z.
func (z *Big) copyAbs(x *Big) *Big {
	if z != x {
		z.precision = x.Precision()
		z.exp = x.exp
		z.compact = x.compact
		if x.IsFinite() && x.isInflated() {
			z.unscaled.Set(&x.unscaled)
		}
	}
	z.form = x.form & ^signbit
	return z
}

// CopySign sets z to x with the sign of y and returns z. It accepts NaN values.
func (z *Big) CopySign(x, y *Big) *Big {
	if debug {
		x.validate()
		y.validate()
	}
	// Pre-emptively capture signbit in case z == y.
	sign := y.form & signbit
	z.copyAbs(x)
	z.form |= sign
	return z
}

// Float64 returns x as a float64 and a bool indicating whether x can be exactly
// represented as a float64. Special values are considered exact; however, special
// values that occur because the magnitude of x is too large to be represented
// as a float64 are not.
func (x *Big) Float64() (f float64, exact bool) {
	if debug {
		x.validate()
	}

	if !x.IsFinite() {
		switch x.form {
		case pinf, ninf:
			return math.Inf(int(x.form & signbit)), true
		case snan, qnan:
			return math.NaN(), true
		case ssnan, sqnan:
			return math.Copysign(math.NaN(), -1), true
		}
	}

	const maxPow10 = 22         // largest exact power of 10
	const maxMantissa = 1 << 52 // largest exact mantissa
	switch {
	case !x.isCompact():
		f, _ = strconv.ParseFloat(x.String(), 64)
		exact = false
		//f, exact = math.Inf(0), false
	case x.compact == 0:
		exact = true
	case x.exp == 0:
		f, exact = float64(x.compact), true
	case x.exp > 0:
		f = float64(x.compact) * math.Pow10(x.exp)
		exact = x.compact < maxMantissa && x.exp < maxPow10
	case x.exp < 0:
		f = float64(x.compact) / math.Pow10(-x.exp)
		exact = x.compact < maxMantissa && x.exp > -maxPow10
	}
	if x.form&signbit != 0 {
		f = math.Copysign(f, -1)
	}
	return f, exact
}

// Float sets z to x and returns z. z is allowed to be nil. The result is
// undefined if z is a NaN value.
func (x *Big) Float(z *big.Float) *big.Float {
	if debug {
		x.validate()
	}

	if z == nil {
		z = new(big.Float)
	}

	switch x.form {
	case finite, finite | signbit:
		if x.compact == 0 {
			z.SetUint64(0)
		} else {
			z.SetRat(x.Rat(nil))
		}
	case pinf, ninf:
		z.SetInf(x.form == pinf)
	default: // snan, qnan, ssnan, sqnan:
		z.SetUint64(0)
	}
	return z
}

// Format implements the fmt.Formatter interface. The following verbs are
// supported:
//
// 	%s: -dddd.dd or -d.dddd±edd, depending on x
// 	%d: same as %s
// 	%v: same as %s
// 	%e: -d.dddd±edd
// 	%E: -d.dddd±Edd
// 	%f: -dddd.dd
// 	%g: same as %f
//
// While width is honored in the same manner as the fmt package (the minimum
// width of the formatted number), precision is the number of significant digits
// in the decimal number. Given %f, however, precision is the number of digits
// following the radix.
//
// Format honors all flags (such as '+' and ' ') in the same manner as the fmt
// package, except for '#'. Unless used in conjunction with %v, %q, or %p, the
// '#' flag will be ignored; decimals have no defined hexadeximal or octal
// representation.
//
// %+v, %#v, %T, %#p, and %p all honor the formats specified in the fmt
// package's documentation.
func (x *Big) Format(s fmt.State, c rune) {
	if debug {
		x.validate()
	}

	prec, hasPrec := s.Precision()
	if !hasPrec {
		prec = x.Precision()
	}
	width, hasWidth := s.Width()
	if !hasWidth {
		width = noWidth
	}

	var (
		hash    = s.Flag('#')
		dash    = s.Flag('-')
		lpZero  = s.Flag('0')
		lpSpace = width != noWidth && !dash && !lpZero
		plus    = s.Flag('+')
		space   = s.Flag(' ')
		f       = formatter{prec: prec, width: width}
		e       = sciE[x.Context.OperatingMode]
	)

	// If we need to left pad then we need to first write our string into an
	// empty buffer.
	tmpbuf := lpZero || lpSpace
	if tmpbuf {
		b := new(strings.Builder)
		b.Grow(x.Precision())
		f.w = b
	} else {
		f.w = stateWrapper{s}
	}

	if plus {
		f.sign = '+'
	} else if space {
		f.sign = ' '
	}

	// noE is a placeholder for formats that do not use scientific notation
	// and don't require 'e' or 'E'
	const noE = 0
	switch c {
	case 's', 'd':
		f.format(x, normal, e)
	case 'q':
		// The fmt package's docs specify that the '+' flag
		// "guarantee[s] ASCII-only output for %q (%+q)"
		f.sign = 0

		// Since no other escaping is needed we can do it ourselves and save
		// whatever overhead running it through fmt.Fprintf would incur.
		quote := byte('"')
		if hash {
			quote = '`'
		}
		f.WriteByte(quote)
		f.format(x, normal, e)
		f.WriteByte(quote)
	case 'e', 'E':
		f.format(x, sci, byte(c))
	case 'f', 'F':
		if !hasPrec {
			prec = 0
		} else {

			// %f's precision means "number of digits after the radix"
			if x.exp > 0 {
				f.prec += x.Precision()
			} else {
				if adj := x.exp + x.Precision(); adj > -f.prec {
					f.prec += adj
				} else {
					f.prec = -f.prec
				}
			}
		}

		f.format(x, plain, noE)
	case 'g', 'G':
		// %g's precision means "number of significant digits"
		f.format(x, plain, noE)

	// Make sure we return from the following two cases.
	case 'v':
		// %v == %s
		if !hash && !plus {
			f.format(x, normal, e)
			break
		}

		// This is the easiest way of doing it. Note we can't use type Big Big,
		// even though it's declared inside a function. Go thinks it's recursive.
		// At least the fields are checked at compile time.
		type Big struct {
			Context   Context
			unscaled  big.Int
			compact   uint64
			exp       int
			precision int
			form      form
		}
		specs := ""
		if dash {
			specs += "-"
		} else if lpZero {
			specs += "0"
		}
		if hash {
			specs += "#"
		} else if plus {
			specs += "+"
		} else if space {
			specs += " "
		}
		fmt.Fprintf(s, "%"+specs+"v", (*Big)(x))
		return
	default:
		fmt.Fprintf(s, "%%!%c(*decimal.Big=%s)", c, x.String())
		return
	}

	// Need padding out to width.
	if f.n < int64(width) {
		switch pad := int64(width) - f.n; {
		case dash:
			io.CopyN(s, spaceReader{}, pad)
		case lpZero:
			io.CopyN(s, zeroReader{}, pad)
		case lpSpace:
			io.CopyN(s, spaceReader{}, pad)
		}
	}

	if tmpbuf {
		// fmt's internal state type implements stringWriter I think.
		io.WriteString(s, f.w.(*strings.Builder).String())
	}
}

var _ fmt.Formatter = (*Big)(nil)

// FMA sets z to (x * y) + u without any intermediate rounding.
func (z *Big) FMA(x, y, u *Big) *Big { return z.Context.FMA(z, x, y, u) }

// Int sets z to x, truncating the fractional portion (if any) and returns z. z
// is allowed to be nil. If x is an infinity or a NaN value the result is
// undefined.
func (x *Big) Int(z *big.Int) *big.Int {
	if debug {
		x.validate()
	}

	if z == nil {
		z = new(big.Int)
	}

	if !x.IsFinite() {
		return z
	}

	if x.isCompact() {
		z.SetUint64(x.compact)
	} else {
		z.Set(&x.unscaled)
	}
	if x.Signbit() {
		z.Neg(z)
	}
	if x.exp == 0 {
		return z
	}
	return bigScalex(z, z, x.exp)
}

// Int64 returns x as an int64, truncating towards zero. The returned boolean
// indicates whether the conversion to an int64 was successful.
func (x *Big) Int64() (int64, bool) {
	if debug {
		x.validate()
	}

	if !x.IsFinite() {
		return 0, false
	}

	// x might be too large to fit into an int64 *now*, but rescaling x might
	// shrink it enough. See issue #20.
	if !x.isCompact() {
		xb := x.Int(nil)
		return xb.Int64(), xb.IsInt64()
	}

	u := x.compact
	if x.exp != 0 {
		var ok bool
		if u, ok = scalex(u, x.exp); !ok {
			return 0, false
		}
	}
	su := int64(u)
	if su >= 0 || x.Signbit() && su == -su {
		if x.Signbit() {
			su = -su
		}
		return su, true
	}
	return 0, false
}

// Uint64 returns x as an int64, truncating towards zero. The returned boolean
// indicates whether the conversion to a uint64 was successful.
func (x *Big) Uint64() (uint64, bool) {
	if debug {
		x.validate()
	}

	if !x.IsFinite() || x.Signbit() {
		return 0, false
	}

	// x might be too large to fit into an uint64 *now*, but rescaling x might
	// shrink it enough. See issue #20.
	if !x.isCompact() {
		xb := x.Int(nil)
		return xb.Uint64(), xb.IsUint64()
	}

	b := x.compact
	if x.exp == 0 {
		return b, true
	}
	return scalex(b, x.exp)
}

// IsFinite returns true if x is finite.
func (x *Big) IsFinite() bool { return x.form & ^signbit == 0 }

// IsNormal returns true if x is normal.
func (x *Big) IsNormal() bool {
	return x.IsFinite() && x.adjusted() >= x.Context.minScale()
}

// IsSubnormal returns true if x is subnormal.
func (x *Big) IsSubnormal() bool {
	return x.IsFinite() && x.adjusted() < x.Context.minScale()
}

// IsInf returns true if x is an infinity according to sign.
// If sign >  0, IsInf reports whether x is positive infinity.
// If sign <  0, IsInf reports whether x is negative infinity.
// If sign == 0, IsInf reports whether x is either infinity.
func (x *Big) IsInf(sign int) bool {
	return sign >= 0 && x.form == pinf || sign <= 0 && x.form == ninf
}

// IsNaN returns true if x is NaN.
// If sign >  0, IsNaN reports whether x is quiet NaN.
// If sign <  0, IsNaN reports whether x is signaling NaN.
// If sign == 0, IsNaN reports whether x is either NaN.
func (x *Big) IsNaN(quiet int) bool {
	return quiet >= 0 && x.form&qnan == qnan || quiet <= 0 && x.form&snan == snan
}

// IsInt reports whether x is an integer. Infinity and NaN values are not
// integers.
func (x *Big) IsInt() bool {
	if debug {
		x.validate()
	}

	if !x.IsFinite() {
		return false
	}

	// 0, 5000, 40
	if x.compact == 0 || x.exp >= 0 {
		return true
	}

	xp := x.Precision()
	exp := x.exp

	// 0.001
	// 0.5
	if -exp >= xp {
		return false
	}

	// 44.00
	// 1.000
	if x.isCompact() {
		for v := x.compact; v%10 == 0; v /= 10 {
			exp++
		}
		// Avoid the overhead of copying x.unscaled if we know for a fact it's not
		// an integer.
	} else if x.unscaled.Bit(0) == 0 {
		v := new(big.Int).Set(&x.unscaled)
		r := new(big.Int)
		for {
			v.QuoRem(v, c.TenInt, r)
			if r.Sign() != 0 {
				break
			}
			exp++
		}
	}
	return exp >= 0
}

// MarshalText implements encoding.TextMarshaler.
func (x *Big) MarshalText() ([]byte, error) {
	if debug {
		x.validate()
	}
	var (
		b = new(bytes.Buffer)
		f = formatter{w: b, prec: x.Precision(), width: noWidth}
		e = sciE[x.Context.OperatingMode]
	)
	b.Grow(x.Precision())
	f.format(x, normal, e)
	return b.Bytes(), nil
}

// Mul sets z to x * y and returns z.
func (z *Big) Mul(x, y *Big) *Big { return z.Context.Mul(z, x, y) }

// Neg sets z to -x and returns z. If x is positive infinity, z will be set to
// negative infinity and visa versa. If x == 0, z will be set to zero as well.
// NaN will result in an error.
func (z *Big) Neg(x *Big) *Big {
	if debug {
		x.validate()
	}
	if !z.invalidContext(z.Context) && !z.checkNaNs(x, x, negation) {
		xform := x.form // copy in case z == x
		z.copyAbs(x)
		if !z.IsFinite() || z.compact != 0 || z.Context.RoundingMode == ToNegativeInf {
			z.form = xform ^ signbit
		}
	}
	return z.Context.round(z)
}

// New creates a new Big decimal with the given value and scale. For example:
//
//  New(1234, 3) // 1.234
//  New(42, 0)   // 42
//  New(4321, 5) // 0.04321
//  New(-1, 0)   // -1
//  New(3, -10)  // 30 000 000 000
//
func New(value int64, scale int) *Big {
	return new(Big).SetMantScale(value, scale)
}

// Payload returns the payload of x, provided x is a NaN value. If x is not a
// NaN value, the result is undefined.
func (x *Big) Payload() Payload {
	if !x.IsNaN(0) {
		return 0
	}
	return Payload(x.compact)
}

// Precision returns the precision of x. That is, it returns the number of
// digits in the unscaled form of x. x == 0 has a precision of 1. The result is
// undefined if x is not finite.
func (x *Big) Precision() int {
	// Cannot call validate since validate calls this method.
	if !x.IsFinite() {
		return 0
	}
	if x.precision == 0 {
		return 1
	}
	return x.precision
}

// Quantize sets z to the number equal in value and sign to z with the scale, n.
func (z *Big) Quantize(n int) *Big { return z.Context.Quantize(z, n) }

// Quo sets z to x / y and returns z.
func (z *Big) Quo(x, y *Big) *Big { return z.Context.Quo(z, x, y) }

// QuoInt sets z to x / y with the remainder truncated. See QuoRem for more
// details.
func (z *Big) QuoInt(x, y *Big) *Big { return z.Context.QuoInt(z, x, y) }

// QuoRem sets z to the quotient x / y and r to the remainder x % y, such that
// x = z * y + r, and returns the pair (z, r).
func (z *Big) QuoRem(x, y, r *Big) (*Big, *Big) {
	return z.Context.QuoRem(z, x, r, y)
}

// Rat sets z to x returns z. z is allowed to be nil. The result is undefined if
// x is an infinity or NaN value.
func (x *Big) Rat(z *big.Rat) *big.Rat {
	if debug {
		x.validate()
	}

	if z == nil {
		z = new(big.Rat)
	}

	if !x.IsFinite() {
		return z.SetInt64(0)
	}

	// Fast path for decimals <= math.MaxInt64.
	if x.IsInt() {
		if u, ok := x.Int64(); ok {
			// If profiled we can call scalex ourselves and save the overhead of
			// calling Int64. But I doubt it'll matter much.
			return z.SetInt64(u)
		}
	}

	num := new(big.Int)
	if x.isCompact() {
		num.SetUint64(x.compact)
	} else {
		num.Set(&x.unscaled)
	}
	if x.exp > 0 {
		checked.MulBigPow10(num, num, uint64(x.exp))
	}
	if x.Signbit() {
		num.Neg(num)
	}

	denom := c.OneInt
	if x.exp < 0 {
		denom = new(big.Int)
		if shift, ok := arith.Pow10(uint64(-x.exp)); ok {
			denom.SetUint64(shift)
		} else {
			denom.Set(arith.BigPow10(uint64(-x.exp)))
		}
	}
	return z.SetFrac(num, denom)
}

// Raw directly returns x's raw compact and unscaled values. Caveat emptor:
// Neither are guaranteed to be valid. Raw is intended to support missing
// functionality outside this package and generally should be avoided.
// Additionally, Raw is the only part of this package's API which is not
// guaranteed to remain stable. This means the function could change or
// disappear at any time, even across minor version numbers.
func Raw(x *Big) (*uint64, *big.Int) { return &x.compact, &x.unscaled }

// Reduce reduces a finite z to its most simplest form.
func (z *Big) Reduce() *Big { return z.Context.Reduce(z) }

// Rem sets z to the remainder x % y. See QuoRem for more details.
func (z *Big) Rem(x, y *Big) *Big { return z.Context.Rem(z, x, y) }

// Round rounds z down to n digits of precision and returns z. The result is
// undefined if z is not finite. No rounding will occur if n <= 0. The result of
// Round will always be within the interval [⌊10**x⌋, z] where x = the precision
// of z.
func (z *Big) Round(n int) *Big {
	ctx := z.Context
	ctx.Precision = n
	return ctx.Round(z)
}

// RoundToInt rounds z down to an integral value.
func (z *Big) RoundToInt() *Big { return z.Context.RoundToInt(z) }

// Scale returns x's scale.
func (x *Big) Scale() int { return -x.exp }

// Scan implements fmt.Scanner.
func (z *Big) Scan(state fmt.ScanState, verb rune) error {
	return z.scan(byteReader{state})
}

var _ fmt.Scanner = (*Big)(nil)

// Set sets z to x and returns z. The result might be rounded depending on z's
// Context, and even if z == x.
func (z *Big) Set(x *Big) *Big { return z.Context.round(z.Copy(x)) }

// setShared sets z to x, but does not copy—z may possibly alias x.
func (z *Big) setShared(x *Big) *Big {
	if debug {
		x.validate()
	}

	if z != x {
		z.precision = x.Precision()
		z.compact = x.compact
		z.form = x.form
		z.exp = x.exp
		z.unscaled = x.unscaled
	}
	return z
}

// SetBigMantScale sets z to the given value and scale.
func (z *Big) SetBigMantScale(value *big.Int, scale int) *Big {
	// Do this first in case value == z.unscaled. Don't want to clobber the sign.
	z.form = finite
	if value.Sign() < 0 {
		z.form |= signbit
	}

	z.unscaled.Abs(value)
	z.compact = c.Inflated
	z.precision = arith.BigLength(value)

	if z.unscaled.IsUint64() {
		if v := z.unscaled.Uint64(); v != c.Inflated {
			z.compact = v
		}
	}

	z.exp = -scale
	return z
}

// SetFloat sets z to x and returns z.
func (z *Big) SetFloat(x *big.Float) *Big {
	if x.IsInf() {
		if x.Signbit() {
			z.form = ninf
		} else {
			z.form = pinf
		}
		return z
	}

	neg := x.Signbit()
	if x.Sign() == 0 {
		if neg {
			z.form |= signbit
		}
		z.compact = 0
		z.precision = 1
		return z
	}

	z.exp = 0
	x0 := new(big.Float).Copy(x).SetPrec(big.MaxPrec)
	x0.Abs(x0)
	if !x.IsInt() {
		for !x0.IsInt() {
			x0.Mul(x0, c.TenFloat)
			z.exp--
		}
	}

	if mant, acc := x0.Uint64(); acc == big.Exact {
		z.compact = mant
		z.precision = arith.Length(mant)
	} else {
		z.compact = c.Inflated
		x0.Int(&z.unscaled)
		z.precision = arith.BigLength(&z.unscaled)
	}
	z.form = finite
	if neg {
		z.form |= signbit
	}
	return z
}

// SetFloat64 sets z to exactly x.
func (z *Big) SetFloat64(x float64) *Big {
	if x == 0 {
		var sign form
		if math.Signbit(x) {
			sign = signbit
		}
		return z.setZero(sign, 0)
	}
	if math.IsNaN(x) {
		var sign form
		if math.Signbit(x) {
			sign = signbit
		}
		return z.setNaN(0, qnan|sign, 0)
	}
	if math.IsInf(x, 0) {
		if math.IsInf(x, 1) {
			z.form = pinf
		} else {
			z.form = ninf
		}
		return z
	}

	// The gist of the following is lifted from math/big/rat.go, but adapted for
	// base-10 decimals.

	const expMask = 1<<11 - 1
	bits := math.Float64bits(x)
	mantissa := bits & (1<<52 - 1)
	exp := int((bits >> 52) & expMask)
	if exp == 0 { // denormal
		exp -= 1022
	} else { // normal
		mantissa |= 1 << 52
		exp -= 1023
	}

	if mantissa == 0 {
		return z.SetUint64(0)
	}

	shift := 52 - exp
	for mantissa&1 == 0 {
		mantissa >>= 1
		shift--
	}

	z.exp = 0
	z.form = finite | form(bits>>63)

	if shift > 0 {
		z.unscaled.SetUint64(uint64(shift))
		z.unscaled.Exp(c.FiveInt, &z.unscaled, nil)
		arith.MulUint64(&z.unscaled, &z.unscaled, mantissa)
		z.exp = -shift
		return z.norm()
	}

	if s := uint(-shift); s < 64 {
		z.compact = mantissa << s
		z.precision = arith.Length(z.compact)
	} else {
		z.compact = c.Inflated
		z.unscaled.SetUint64(mantissa)
		z.unscaled.Lsh(&z.unscaled, s)
		z.norm()
	}
	return z
}

// SetInf sets z to -Inf if signbit is set or +Inf is signbit is not set, and
// returns z.
func (z *Big) SetInf(signbit bool) *Big {
	if signbit {
		z.form = ninf
	} else {
		z.form = pinf
	}
	return z
}

// SetMantScale sets z to the given value and scale.
func (z *Big) SetMantScale(value int64, scale int) *Big {
	z.SetUint64(arith.Abs(value))
	z.exp = -scale // compiler should optimize out z.exp = 0 in SetUint64
	if value < 0 {
		z.form |= signbit
	}
	return z
}

// setNaN is an internal NaN-setting method that panics when the OperatingMode
// is Go.
func (z *Big) setNaN(c Condition, f form, p Payload) *Big {
	z.form = f
	z.compact = uint64(p)
	z.Context.Conditions |= c
	if z.Context.OperatingMode == Go {
		panic(ErrNaN{Msg: z.Context.Conditions.String()})
	}
	return z
}

// SetNaN sets z to a signaling NaN if signal is true or quiet NaN otherwise and
// returns z. No conditions are raised.
func (z *Big) SetNaN(signal bool) *Big {
	if signal {
		z.form = snan
	} else {
		z.form = qnan
	}
	z.compact = 0 // payload
	return z
}

// SetRat sets z to to the possibly rounded value of x and return z.
func (z *Big) SetRat(x *big.Rat) *Big {
	if x.IsInt() {
		return z.Context.round(z.SetBigMantScale(x.Num(), 0))
	}
	var num, denom Big
	num.SetBigMantScale(x.Num(), 0)
	denom.SetBigMantScale(x.Denom(), 0)
	return z.Quo(&num, &denom)
}

// SetScale sets z's scale to scale and returns z.
func (z *Big) SetScale(scale int) *Big {
	z.exp = -scale
	return z
}

// Regexp matches any valid string representing a decimal that can be passed to
// SetString.
var Regexp = regexp.MustCompile(`(?i)(([+-]?(\d+\.\d*|\.?\d+)([eE][+-]?\d+)?)|(inf(infinity)?))|([+-]?([sq]?nan\d*))`)

// SetString sets z to the value of s, returning z and a bool indicating success.
// s must be a string in one of the following formats:
//
// 	1.234
// 	1234
// 	1.234e+5
// 	1.234E-5
// 	0.000001234
// 	Inf
// 	NaN
// 	qNaN
// 	sNaN
//
// Each value may be preceded by an optional sign, ``-'' or ``+''. ``Inf'' and
// ``NaN'' map to ``+Inf'' and ``qNaN'', respectively. NaN values may have
// optional diagnostic information, represented as trailing digits; for example,
// ``NaN123''. These digits are otherwise ignored but are included for
// robustness.
func (z *Big) SetString(s string) (*Big, bool) {
	if err := z.scan(strings.NewReader(s)); err != nil {
		return nil, false
	}
	return z, true
}

func (z *Big) setTriple(compact uint64, sign form, exp int) *Big {
	z.compact = compact
	z.precision = arith.Length(compact)
	z.exp = exp
	z.form = finite | sign
	return z
}

func (z *Big) setZero(sign form, exp int) *Big {
	z.compact = 0
	z.precision = 1
	z.exp = exp
	z.form = finite | sign
	return z
}

// SetUint64 is shorthand for SetMantScale(x, 0) for an unsigned integer.
func (z *Big) SetUint64(x uint64) *Big {
	z.compact = x
	z.precision = arith.Length(x)
	z.exp = 0
	z.form = finite
	return z
}

// ord returns similar to Sign except -Inf is -2 and +Inf is +2.
func (x *Big) ord(abs bool) int {
	if x.form&inf != 0 {
		if x.form == pinf || abs {
			return +2
		}
		return -2
	}
	r := x.Sign()
	if abs && r < 0 {
		r = -r
	}
	return r
}

// Sign returns:
//
//	-1 if x <  0
//	 0 if x == 0
//	+1 if x >  0
//
// The result is undefined if x is a NaN value.
func (x *Big) Sign() int {
	if debug {
		x.validate()
	}

	if (x.IsFinite() && x.compact == 0) || x.IsNaN(0) {
		return 0
	}
	if x.form&signbit != 0 {
		return -1
	}
	return 1
}

// Signbit returns true if x is negative, negative infinity, negative zero, or
// negative NaN.
func (x *Big) Signbit() bool {
	if debug {
		x.validate()
	}
	return x.form&signbit != 0
}

// String returns the string representation of x. It's equivalent to the %s verb
// discussed in the Format method's documentation. Special cases depend on the
// OperatingMode.
func (x *Big) String() string {
	var (
		b = new(strings.Builder)
		f = formatter{w: b, prec: x.Precision(), width: noWidth}
		e = sciE[x.Context.OperatingMode]
	)
	b.Grow(x.Precision())
	f.format(x, normal, e)
	return b.String()
}

var _ fmt.Stringer = (*Big)(nil)

// Sub sets z to x - y and returns z.
func (z *Big) Sub(x, y *Big) *Big { return z.Context.Sub(z, x, y) }

// UnmarshalText implements encoding.TextUnmarshaler.
func (z *Big) UnmarshalText(data []byte) error {
	return z.scan(bytes.NewReader(data))
}

var _ encoding.TextUnmarshaler = (*Big)(nil)

// validate ensures x's internal state is correct. There's no need for it to
// have good performance since it's for debug == true only.
func (x *Big) validate() {
	defer func() {
		if err := recover(); err != nil {
			pc, _, _, ok := runtime.Caller(4)
			if caller := runtime.FuncForPC(pc); ok && caller != nil {
				fmt.Println("called by:", caller.Name())
			}
			type Big struct {
				Context   Context
				unscaled  big.Int
				compact   uint64
				exp       int
				precision int
				form      form
			}
			fmt.Printf("%#v\n", (*Big)(x))
			panic(err)
		}
	}()
	switch x.form {
	case finite, finite | signbit:
		if x.isInflated() {
			if x.unscaled.IsUint64() && x.unscaled.Uint64() != c.Inflated {
				panic(fmt.Sprintf("inflated but unscaled == %d", x.unscaled.Uint64()))
			}
			if x.unscaled.Sign() < 0 {
				panic("x.unscaled.Sign() < 0")
			}
			if bl, xp := arith.BigLength(&x.unscaled), x.precision; bl != xp {
				panic(fmt.Sprintf("BigLength (%d) != x.Precision (%d)", bl, xp))
			}
		}
		if x.isCompact() {
			if bl, xp := arith.Length(x.compact), x.Precision(); bl != xp {
				panic(fmt.Sprintf("BigLength (%d) != x.Precision() (%d)", bl, xp))
			}
		}
	case snan, ssnan, qnan, sqnan, pinf, ninf:
		// OK
	case nan:
		panic(x.form.String())
	default:
		panic(fmt.Sprintf("invalid form %s", x.form))
	}
}
