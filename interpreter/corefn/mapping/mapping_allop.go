package mapping

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaFunction {
	return []shared.TaFunction{
		Map,
	}
}
