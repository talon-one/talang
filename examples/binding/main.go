package main

import (
	"fmt"

	"github.com/talon-one/talang"
	"github.com/talon-one/talang/token"
)

func main() {
	interp := talang.MustNewInterpreter()
	// Set the Name to Joe
	interp.Set("Name", token.NewString("Joe"))
	result, err := interp.LexAndEvaluate(`(+ "Hello " (. Name) "!")`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Stringify())

	// You can also set the binding in an expression
	_, err = interp.LexAndEvaluate(`(set Name "Doe")`)
	if err != nil {
		panic(err)
	}

	result, err = interp.LexAndEvaluate(`(+ "Hello " (. Name) "!")`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Stringify())
}
