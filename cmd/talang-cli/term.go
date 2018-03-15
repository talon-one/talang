package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var term *readline.Instance

var prompt = color.HiCyanString("talang") + color.WhiteString("> ")

func beginTerm() {
	var completer = readline.NewPrefixCompleter(
		readline.PcItemDynamic(loadHints()),
	)

	// Setup readline
	newTerm, err := readline.NewEx(&readline.Config{
		Prompt:            prompt,
		HistoryFile:       getUserHistoryFilePath(),
		AutoComplete:      completer,
		InterruptPrompt:   "^C",
		EOFPrompt:         "^D",
		HistorySearchFold: true,
	})
	if err != nil {
		panic(err)
	}

	term = newTerm
	printOut("Welcome to talang cli! Enter :help to get some help, :exit to exit.")
}

func endTerm() {
	term.Close()
}

func out() (o io.Writer) {
	o = os.Stdout
	if term != nil {
		o = term.Stdout()
	}
	return o
}

func printOut(msg string, args ...interface{}) {
	fmt.Fprintf(out(), msg+"\n", args...)
}

func printResult(msg string) {
	fmt.Fprintln(out(), color.GreenString(msg))
}

func printErr(e error) {
	fmt.Fprintln(out(), color.RedString(e.Error()))
}

// User history file helpers
func getUserHistoryFilePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	var path = filepath.Join(usr.HomeDir, ".tex")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	path = filepath.Join(path, "history.tmp")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		var file, _ = os.Create(path)
		defer file.Close()
	}

	return path
}

// todo: this can be made better
func loadHints() func(string) []string {
	return func(line string) []string {
		return allHints
	}
}
