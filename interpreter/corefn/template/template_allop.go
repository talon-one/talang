package template

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaSignature {
	return []shared.TaSignature{
		GetTemplate,
		SetTemplate,
	}
}
