package main

import (
	"fmt"

	"github.com/talon-one/talang"
)

func main() {
	interp := talang.MustNewInterpreter()
	result, err := interp.LexAndEvaluate(`(+ 1 2)`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Stringify())
}
