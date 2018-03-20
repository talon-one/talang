package main

import (
	"io"
	"runtime"
	"strings"

	"github.com/chzyer/readline"
	"github.com/talon-one/talang/interpreter"

	"github.com/talon-one/talang"
	"gopkg.in/alecthomas/kingpin.v2"
)

var interp *talang.Interpreter

var interpFunctions []interpreter.TaFunction
var fnHints []string
var cmdHints []string
var allHints []string

var (
	commandArg = kingpin.Arg("command", "command to run directly").String()
)

func main() {
	kingpin.Version("1.0.0")
	kingpin.Parse()

	createCommands()

	interp = talang.MustNewInterpreter()
	interpFunctions = interp.AllFunctions()

	// filter out functions with the same name

	for _, fn := range interpFunctions {
		addFn := true
		for _, hf := range fnHints {
			if fn.Name == hf {
				addFn = false
				break
			}
		}
		if addFn {
			fnHints = append(fnHints, "("+fn.Name)
			cmdHints = append(cmdHints, ":fn "+fn.Name)
		}
	}

	for command := range commands {
		if command != ":fn" {
			cmdHints = append(cmdHints, command)
		}
	}

	allHints = append(fnHints, cmdHints...)

	defer func() {
		if err := recover(); err != nil {
			if _, ok := err.(runtime.Error); ok {
				panic(err)
			}
			printErr(err.(error))
		}
	}()

	if commandArg != nil && len(*commandArg) > 0 {
		evaluateLine(*commandArg)
		return
	}

	beginTerm()
	defer endTerm()

	// Loop for user input
	for {
		line, err := term.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) > 0 {
			evaluateLine(line)
		}
	}
}

func evaluateLine(line string) {
	if command, ok := isPromptCommand(line); ok {
		if err := runCommand(command); err != nil {
			printOut(err.Error())
		}
	} else {
		parsed, err := talang.Parse(line)
		if err != nil {
			printErr(err)
		} else {
			if err := interp.Evaluate(parsed); err != nil {
				printErr(err)
			} else {
				printResult(parsed.Stringify())
			}
		}
	}
}
