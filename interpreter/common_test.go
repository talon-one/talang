package interpreter

import (
	"log"
	"os"
)

func mustNewInterpreterWithLogger() *Interpreter {
	interp := MustNewInterpreter()
	interp.Logger = log.New(os.Stdout, "", log.LstdFlags)
	return interp
}
