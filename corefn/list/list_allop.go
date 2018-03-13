package list

import "github.com/talon-one/talang/interpreter"

func AllOperations() []interpreter.TaFunction {
	return []interpreter.TaFunction{
		List,
		Head,
		Tail,
		Drop,
		Item,
		Push,
		Map,
		Sort,
		Min,
		Max,
		Append,
		Count,
		Reverse,
		Join,
	}
}
