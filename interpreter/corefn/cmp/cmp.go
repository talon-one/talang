package cmp

import (
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/talon-one/talang/interpreter/internal"
	"github.com/talon-one/talang/term"
)

func Equal(interp *internal.Interpreter, args ...term.Term) (string, error) {
	argc := len(args)
	if argc < 2 {
		return "", errors.New("invalid or missing arguments")
	}

	for i := 1; i < argc; i++ {
		if args[0].Text != args[i].Text {
			return "false", nil
		}
	}
	return "true", nil
}

func NotEqual(interp *internal.Interpreter, args ...term.Term) (string, error) {
	argc := len(args)
	if argc < 2 {
		return "", errors.New("invalid or missing arguments")
	}

	for i := 1; i < argc; i++ {
		if args[0].Text == args[i].Text {
			return "false", nil
		}
	}
	return "true", nil
}

func GreaterThan(interp *internal.Interpreter, args ...term.Term) (string, error) {
	argc := len(args)
	if argc < 2 {
		return "", errors.New("invalid or missing arguments")
	}

	var d decimal.Decimal
	for i := 0; i < argc; i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if i == 0 {
			d = args[i].Decimal
		} else {
			if !d.GreaterThan(args[i].Decimal) {
				return "false", nil
			}
		}
	}
	return "true", nil
}

func LessThan(interp *internal.Interpreter, args ...term.Term) (string, error) {
	argc := len(args)
	if argc < 2 {
		return "", errors.New("invalid or missing arguments")
	}

	var d decimal.Decimal
	for i := 0; i < argc; i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if i == 0 {
			d = args[i].Decimal
		} else {
			if !d.LessThan(args[i].Decimal) {
				return "false", nil
			}
		}
	}
	return "true", nil
}

func GreaterThanOrEqual(interp *internal.Interpreter, args ...term.Term) (string, error) {
	argc := len(args)
	if argc < 2 {
		return "", errors.New("invalid or missing arguments")
	}

	var d decimal.Decimal
	for i := 0; i < argc; i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if i == 0 {
			d = args[i].Decimal
		} else {
			if !d.GreaterThanOrEqual(args[i].Decimal) {
				return "false", nil
			}
		}
	}
	return "true", nil
}

func LessThanOrEqual(interp *internal.Interpreter, args ...term.Term) (string, error) {
	argc := len(args)
	if argc < 2 {
		return "", errors.New("invalid or missing arguments")
	}

	var d decimal.Decimal
	for i := 0; i < argc; i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if i == 0 {
			d = args[i].Decimal
		} else {
			if !d.LessThanOrEqual(args[i].Decimal) {
				return "false", nil
			}
		}
	}
	return "true", nil
}

func Between(interp *internal.Interpreter, args ...term.Term) (string, error) {
	argc := len(args)
	if argc < 3 {
		return "", errors.New("invalid or missing arguments")
	}

	min := args[argc-2]
	max := args[argc-1]
	if !min.IsDecimal() {
		return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", min.Text)
	}
	if !max.IsDecimal() {
		return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", max.Text)
	}

	argc -= 2

	for i := 0; i < argc; i++ {
		if !args[i].IsDecimal() {
			return "", errors.Errorf("invalid or missing arguments: `%s' is not a decimal", args[i].Text)
		}
		if !args[i].Decimal.GreaterThanOrEqual(min.Decimal) || !args[i].Decimal.LessThanOrEqual(max.Decimal) {
			return "false", nil
		}
	}
	return "true", nil
}
