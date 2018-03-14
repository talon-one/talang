package string

import "github.com/talon-one/talang/interpreter"

func AllOperations() []interpreter.TaFunction {
	return []interpreter.TaFunction{
		Add,
		Concat,
		Contains,
		NotContains,
		StartsWith,
		EndsWith,
		Regexp,
		LastName,
	}
}
