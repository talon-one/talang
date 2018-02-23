package string

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaFunction {
	return []shared.TaFunction{
		Add,
		Concat,
		Contains,
		NotContains,
		StartsWith,
		EndsWith,
		Regexp,
	}
}
