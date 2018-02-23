package math

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaFunction {
	return []shared.TaFunction{
		Add,
		Sub,
		Mul,
		Div,
		Mod,
		Floor,
		Ceil,
	}
}
