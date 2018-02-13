package list

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaSignature {
	return []shared.TaSignature{
		List,
		Head,
		Tail,
		Item,
	}
}
