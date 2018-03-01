package mapping

import "github.com/talon-one/talang/interpreter"

func AllOperations() []interpreter.TaFunction {
	return []interpreter.TaFunction{
		Map,
	}
}
