package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/talon-one/talang/interpreter"
)

type commandDef struct {
	function    func([]string)
	description string
}

var commands map[string]*commandDef

var whitespace = regexp.MustCompile(`\s+`)

func parseArgs(input string) []string {
	parts := whitespace.Split(input, -1)

	args := make([]string, 0, len(parts))

	for i := 0; i < len(parts); i++ {
		p := strings.TrimSpace(parts[i])
		if len(p) > 0 {
			args = append(args, p)
		}
	}

	return args
}

func isPromptCommand(l string) (command string, is bool) {
	command = strings.Trim(l, " ")
	return command, 0 == strings.Index(command, ":")
}

func runCommand(call string) error {
	args := parseArgs(call)

	if len(args) <= 0 {
		return errors.New("unknown command: '" + call + "'")
	}
	com, ok := commands[strings.ToLower(args[0])]
	if ok == false {
		return errors.New("unknown command: '" + call + "'")
	}
	com.function(args[1:])
	return nil
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
		argumentList[j] = color.WhiteString(fn.Arguments[j].String())
	}

	arguments := strings.Join(argumentList, ", ")
	if fn.IsVariadic {
		arguments += "..."
	}

	str := fmt.Sprintf("%s(%s)%s\n    %s", color.YellowString(fn.Name), arguments, color.WhiteString(fn.Returns.String()), strings.TrimSpace(fn.Description))
	if examples {
		str += "\n" + color.New(color.FgHiCyan, color.Underline).Sprint("Examples") + "\n" + strings.TrimSpace(fn.Example) + "\n\n"
	}

	return str
}

func createCommands() {
	commands = make(map[string]*commandDef)

	//
	// General cli commands
	//

	commandHelp := commandDef{
		function: func(args []string) {
			var sortedCommands []string
			for comname := range commands {
				sortedCommands = append(sortedCommands, comname)
			}

			sort.Strings(sortedCommands)

			printOut("\nCommands available:\n\n")
			for _, scomname := range sortedCommands {
				printOut("  %-24s%s\n", scomname, commands[scomname].description)
			}
		},
		description: "displays this help message",
	}

	commandQuit := commandDef{
		function: func(args []string) {
			/**
			 * Using the quit command is the proper way of exiting the cli.
			 * If any cleanup is required it should be done here.
			 */
			endTerm()
			printOut("bye bye!\n\n")
			os.Exit(0)
		},
		description: "quits the talang cli",
	}

	commandDebug := commandDef{
		function: func(args []string) {
			if interp.Logger == nil {
				interp.Logger = log.New(out(), "", 0)
				printOut("Debug mode enabled")
			} else {
				interp.Logger = nil
				printOut("Debug mode disabled")
			}
		},
		description: "enable or disable debug messages",
	}

	commandFn := commandDef{
		function: func(args []string) {
			if len(args) > 0 {
				if !fn(out(), strings.Join(args, " ")) {
					printOut(color.RedString("Unable to find a function matching `%s'", args))
				}
			} else {
				printOut(color.RedString("You need to specify a search string"))
			}
		},
		description: "show the function details",
	}

	// Help
	commands[":?"] = &commandHelp
	commands[":help"] = &commandHelp

	commands[":fn"] = &commandFn

	// Quit
	commands[":q"] = &commandQuit
	commands[":quit"] = &commandQuit

	commands[":d"] = &commandDebug
	commands[":debug"] = &commandDebug
}
