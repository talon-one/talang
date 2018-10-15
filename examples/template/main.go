package main

import (
	"fmt"

	"github.com/talon-one/talang"
	"github.com/talon-one/talang/interpreter"
	"github.com/talon-one/talang/token"
)

func main() {
	interp := talang.MustNewInterpreter()
	interp.MustRegisterTemplate(
		interpreter.TaTemplate{
			CommonSignature: interpreter.CommonSignature{
				Name: "PrintName",
				Arguments: []token.Kind{
					token.String | token.Token, // accept strings and tokens (we need tokens to accept bindings)
					token.String | token.Token,
				},
				Returns: token.String,
			},
			Template: *talang.MustLex(`(+ "Hello " (# 0) " " (# 1) "!")`),
		},
	)

	result, err := interp.LexAndEvaluate(`(! PrintName Joe Doe)`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Stringify())

	interp.Set("Name", token.NewString("Joe"))
	interp.Set("Surname", token.NewString("Doe"))

	result, err = interp.LexAndEvaluate(`(! PrintName (. Name) (. Surname))`)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Stringify())
}
