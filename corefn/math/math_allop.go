package math

import "github.com/talon-one/talang/interpreter"

func AllOperations() []interpreter.TaFunction {
	return []interpreter.TaFunction{
		Add,
		Sub,
		Mul,
		Div,
		Mod,
		Floor,
		Ceil,
	}
}
