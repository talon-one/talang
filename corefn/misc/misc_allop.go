package misc

import "github.com/talon-one/talang/interpreter"

func AllOperations() []interpreter.TaFunction {
	return []interpreter.TaFunction{
		Noop,
		ToString,
		Not,
		Catch,
		Do,
		DoLegacy,
		SafeRead,
		Identity,
	}
}
