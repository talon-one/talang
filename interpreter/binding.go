package interpreter

import "github.com/talon-one/talang/block"

type Binding struct {
	Value    *block.Block
	Children map[string]Binding
}
