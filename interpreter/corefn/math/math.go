package math

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/talon-one/talang/interpreter/internal"
	"github.com/talon-one/talang/term"
)

func Add(interp *internal.Interpreter, args ...term.Term) (string, error) {
	if len(args) < 2 {
		return "", errors.New("invalid or missing arguments")
	}
	d := decimal.Zero
	for i := 0; i < len(args); i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		d = d.Add(args[i].Decimal)
	}
	return d.String(), nil
}

func Sub(interp *internal.Interpreter, args ...term.Term) (string, error) {
	if len(args) < 2 {
		return "", errors.New("invalid or missing arguments")
	}
	var d decimal.Decimal
	for i := 0; i < len(args); i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if i == 0 {
			d = args[i].Decimal
		} else {
			d = d.Sub(args[i].Decimal)
		}
	}
	return d.String(), nil
}

func Mul(interp *internal.Interpreter, args ...term.Term) (string, error) {
	if len(args) < 2 {
		return "", errors.New("invalid or missing arguments")
	}
	var d decimal.Decimal
	for i := 0; i < len(args); i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if i == 0 {
			d = args[i].Decimal
		} else {
			d = d.Mul(args[i].Decimal)
		}
	}
	return d.String(), nil
}

func Div(interp *internal.Interpreter, args ...term.Term) (string, error) {
	if len(args) < 2 {
		return "", errors.New("invalid or missing arguments")
	}
	var d decimal.Decimal
	for i := 0; i < len(args); i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if i == 0 {
			d = args[i].Decimal
		} else {
			d = d.Div(args[i].Decimal)
		}
	}
	return d.String(), nil
}

func Mod(interp *internal.Interpreter, args ...term.Term) (string, error) {
	if len(args) < 2 {
		return "", errors.New("invalid or missing arguments")
	}
	var d decimal.Decimal
	for i := 0; i < len(args); i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if i == 0 {
			d = args[i].Decimal
		} else {
			d = d.Mod(args[i].Decimal)
		}
	}
	return d.String(), nil
}

func Ceil(interp *internal.Interpreter, args ...term.Term) (string, error) {
	if len(args) != 1 {
		return "", errors.New("invalid or missing arguments")
	}
	if !args[0].IsDecimal() {
		return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[0].Text)
	}
	return args[0].Decimal.Ceil().String(), nil
}

func Floor(interp *internal.Interpreter, args ...term.Term) (string, error) {
	if len(args) != 1 {
		return "", errors.New("invalid or missing arguments")
	}
	if !args[0].IsDecimal() {
		return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[0].Text)
	}
	return args[0].Decimal.Floor().String(), nil
}
