package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/talon-one/talang/interpreter"

	"regexp"

	"github.com/fatih/color"
	"github.com/talon-one/talang"
	"golang.org/x/crypto/ssh/terminal"
)

var interp *talang.Interpreter

var interpFunctions []interpreter.TaFunction
var oldState *terminal.State

func main() {
	flag.Parse()

	interp = talang.MustNewInterpreter()
	interpFunctions = interp.AllFunctions()

	args := flag.Args()
	if len(args) > 0 {
		cmd := strings.Join(args, " ")
		runCommand(os.Stdout, cmd)
		return
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGABRT,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		<-sigc
		if oldState != nil {
			terminal.Restore(0, oldState)
		}
		os.Exit(1)
	}()

	var err error
	oldState, err = terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}

	prompt := color.CyanString("talang> ")

	t := terminal.NewTerminal(os.Stdout, prompt)

	io.WriteString(t, "Welcome to talang cli enter .help to get some help, .exit to exit.\n")

	for {
		line, err := t.ReadLine()
		if err != nil {
			t.Write([]byte(err.Error()))
			break
		}
		runCommand(t, line)
	}

	terminal.Restore(0, oldState)
}

func runCommand(out io.Writer, s string) {
	if runInternalCommand(out, s) {
		return
	}

	if fn(out, s) {
		return
	}

	result, err := interp.LexAndEvaluate(s)
	if err != nil {
		color.New(color.FgRed).Fprintf(out, "Error: %s\n", err.Error())
		return
	}
	fmt.Fprintf(out, "%s\n", result.Stringify())
}

func runInternalCommand(out io.Writer, s string) bool {
	s = strings.TrimSpace(strings.ToLower(s))

	split := func(s string) []string {
		p := strings.Split(s, " ")
		parts := make([]string, 0, cap(p))
		for i := 0; i < len(p); i++ {
			part := strings.TrimSpace(p[i])
			if len(part) > 0 {
				parts = append(parts, part)
			}
		}
		return parts
	}

	args := split(s)

	if len(args) <= 0 {
		return false
	}

	switch args[0] {
	case ".quit":
		fallthrough
	case ".exit":
		os.Exit(0)
		return true
	case ".help":
		runHelp(out, args[1:])
		return true
	case ".fns":
		runFns(out, args[1:])
		return true
	}

	return false
}

func runHelp(out io.Writer, args []string) {
	fmt.Fprint(out, `Help topics:
.fns                   -        list all functions
.<function name>       -        get help for a function
.fns <function name>   -        get help for a function
`)
}

func runFns(out io.Writer, args []string) {
	if len(args) <= 0 {
		for i := 0; i < len(interpFunctions); i++ {
			fmt.Fprintln(out, printFunction(&interpFunctions[i], false))
		}
		return
	}

	fn(out, args[0])
}

func fn(out io.Writer, s string) bool {
	matched := false
	r, err := regexp.Compile(s)
	if err == nil {
		for i := 0; i < len(interpFunctions); i++ {
			if r.MatchString(interpFunctions[i].Name) {
				fmt.Fprintln(out, printFunction(&interpFunctions[i], true))
				matched = true
			}
		}
	} else {
		for i := 0; i < len(interpFunctions); i++ {
			if strings.EqualFold(interpFunctions[i].Name, s) {
				fmt.Fprintln(out, printFunction(&interpFunctions[i], true))
				matched = true
			}
		}
	}
	return matched
}

func printFunction(fn *interpreter.TaFunction, examples bool) string {
	argumentList := make([]string, len(fn.Arguments))
	for j := 0; j < len(fn.Arguments); j++ {
		argumentList[j] = fn.Arguments[j].String()
		if strings.HasSuffix(argumentList[j], "Kind") {
			argumentList[j] = argumentList[j][:len(argumentList[j])-4]
		}
		argumentList[j] = color.WhiteString(argumentList[j])
	}

	arguments := strings.Join(argumentList, ", ")
	if fn.IsVariadic {
		arguments += "..."
	}

	returns := fn.Returns.String()
	if strings.HasSuffix(returns, "Kind") {
		returns = returns[:len(returns)-4]
	}
	returns = color.WhiteString(returns)

	str := fmt.Sprintf("%s(%s)%s\n    %s", color.YellowString(fn.Name), arguments, returns, strings.TrimSpace(fn.Description))
	if examples {
		str += "\n" + color.New(color.FgHiCyan, color.Underline).Sprint("Examples") + "\n" + strings.TrimSpace(fn.Example) + "\n\n"
	}

	return str
}
