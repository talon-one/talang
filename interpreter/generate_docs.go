// +build -ignore

package main

import (
	"encoding/json"
	"os"
	"sort"
	"strings"

	"github.com/talon-one/talang/interpreter"
)

func main() {
	f, err := os.Create("functions.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	interp := interpreter.MustNewInterpreter()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "    ")
	encoder.SetEscapeHTML(false)

	type fn struct {
		IsVariadic bool
		Arguments  []string
		Name       string
		Returns    string
	}

	fns := make([]fn, len(interp.Functions()))

	for i, f := range interp.Functions() {
		arguments := make([]string, len(f.Arguments))
		for j, a := range f.Arguments {
			arg := a.String()
			if strings.HasSuffix(arg, "Kind") {
				arguments[j] = arg[:len(arg)-4]
			} else {
				arguments[j] = arg
			}
		}

		returns := f.Returns.String()
		if strings.HasSuffix(returns, "Kind") {
			returns = returns[:len(returns)-4]
		}

		fns[i] = fn{
			IsVariadic: f.IsVariadic,
			Name:       f.Name,
			Arguments:  arguments,
			Returns:    returns,
		}
	}

	sort.Slice(fns, func(i, j int) bool {
		return strings.Compare(fns[i].Name, fns[j].Name) < 0
	})
	encoder.Encode(fns)
}
