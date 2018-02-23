package misc

import "github.com/talon-one/talang/interpreter/shared"

func AllOperations() []shared.TaFunction {
	return []shared.TaFunction{
		Misc,
		Noop,
		ToString,
	}
}
