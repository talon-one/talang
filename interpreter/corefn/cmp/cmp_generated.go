package cmp

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaSignature {
	return []shared.TaSignature{
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
