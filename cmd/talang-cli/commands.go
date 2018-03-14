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
	function    func(string)
	description string
}

var commands map[string]*commandDef

var whitespace, _ = regexp.Compile("\\s+")

func firstWord(input string) (string, string) {
	parts := whitespace.Split(strings.Trim(input, " "), 2)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}

func isPromptCommand(l string) (command string, is bool) {
	command = strings.Trim(l, " ")
	return command, 0 == strings.Index(command, ":")
}

func runCommand(call string) error {
	name, args := firstWord(call)
	com := commands[name]
	if com == nil {
		return errors.New("unknown command: '" + call + "'")
	}
	com.function(args)
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

func createCommands() {
	commands = map[string]*commandDef{}

	//
	// General cli commands
	//

	commandHelp := commandDef{
		function: func(args string) {
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
		function: func(args string) {
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
		function: func(args string) {
			first, _ := firstWord(args)
			switch first {
			case "debug":
				if interp.Logger == nil {
					interp.Logger = log.New(term.Stdout(), "", 0)
					printOut("Debug mode enabled")
				} else {
					interp.Logger = nil
					printOut("Debug mode disabled")
				}
			default:
				printOut("unknown command: set %s", first)
			}
		},
		description: "enable or disable debug messages",
	}

	commandFn := commandDef{
		function: func(args string) {
			if len(args) > 0 {
				if !fn(term.Stdout(), args) {
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
