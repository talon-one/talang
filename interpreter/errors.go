package interpreter

import (
	"fmt"
	"strings"

	"github.com/talon-one/talang/token"
)

type MaxRecursiveLevelReachedError struct {
	Level int
}

func (err MaxRecursiveLevelReachedError) Error() string {
	return fmt.Sprintf("Max recursive level (%d) reached", err.Level)
}

type FunctionNotFoundError struct {
	token           *token.TaToken
	CollectedErrors []error
}

func (err FunctionNotFoundError) Error() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "Found no function for `%s'", err.token.Stringify())

	for i := 0; i < len(err.CollectedErrors); i++ {
		builder.WriteRune('\n')
		builder.WriteString("  ")
		builder.WriteString(err.CollectedErrors[i].Error())
	}
	return builder.String()
}

type FunctionError struct {
	function   *TaFunction
	StackTrace error
	error
}

func (err FunctionError) Error() string {
	return fmt.Sprintf("Error in function `%s': %s", err.function.CommonSignature.String(), err.error.Error())
}

type FunctionErrors struct {
	Errors []FunctionError
	token  *token.TaToken
}

func (err FunctionErrors) Error() string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "Errors in `%s':", err.token.Stringify())

	for i := 0; i < len(err.Errors); i++ {
		builder.WriteRune('\n')
		builder.WriteString("  ")
		builder.WriteString(err.Errors[i].Error())
	}
	return builder.String()
}

type FunctionNotRanError struct {
	function *TaFunction
	Reason   error
}

func (err FunctionNotRanError) Error() string {
	return fmt.Sprintf("Not Running Function `%s': %s", err.function.CommonSignature.String(), err.Reason.Error())
}
