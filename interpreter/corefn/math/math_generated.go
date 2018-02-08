package math

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaSignature {
	return []shared.TaSignature{
		Add,
		Sub,
		Mul,
		Div,
		Mod,
		Ceil,
		Floor,
	}
}
