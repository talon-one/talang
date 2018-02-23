package list

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaFunction {
	return []shared.TaFunction{
		List,
		Head,
		Tail,
		Drop,
		Item,
	}
}
