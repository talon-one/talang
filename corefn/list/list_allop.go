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
		MapLegacy,
		Sort,
		Min,
		Max,
		Append,
		Count,
		Reverse,
		Join,
		IsEmpty,
		Split,
		Exists,
		ExistsLegacy,
		Sum,
		Every,
		EveryLegacy,
		SortByNumber,
	}
}
