package cmp

import "github.com/talon-one/talang/interpreter"

func AllOperations() []interpreter.TaFunction {
	return []interpreter.TaFunction{
		Equal,
		NotEqual,
		GreaterThanDecimal,
		GreaterThanTime,
		LessThanDecimal,
		LessThanTime,
		GreaterThanOrEqualDecimal,
		GreaterThanOrEqualTime,
		LessThanOrEqualDecimal,
		LessThanOrEqualTime,
		BetweenDecimal,
		BetweenTime,
		Or,
	}
}
