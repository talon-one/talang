package string

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaSignature {
	return []shared.TaSignature{
		Add,
		Concat,
		Contains,
		NotContains,
		StartsWith,
		EndsWith,
		Regexp,
	}
}
