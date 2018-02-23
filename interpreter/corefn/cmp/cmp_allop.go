package cmp

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaFunction {
	return []shared.TaFunction{
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
	}
}
