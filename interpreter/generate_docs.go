// +build -ignore

package main

import (
	"encoding/json"
	"flag"
	"os"
	"sort"
	"strings"

	"html/template"
	texttemplate "text/template"

	"github.com/talon-one/talang/interpreter"
)

var flagOutput = flag.String("format", "md", "format to use for output")

const htmlTemplate string = `
<html>
	<style type="text/css">
	p {
		font-family: monospace;
		display: block;
		margin: .4rem 0;
		padding: 0;
	}
	</style>
	{{ range $index, $element := . }}
		<p><b>{{ $element.Name -}}</b>
			(
			{{- range $i,$arg := $element.Arguments -}}
				{{if $i}}, {{end}}{{ $arg -}}
			{{ end -}}
			{{ if $element.IsVariadic }}...{{ end -}}
			)
			{{- $element.Returns }}
		</p>
	{{ end }}
</html>
`

const markdownTemplate string = `
# Functions
{{ range $index, $element := . }}
    {{ $element.Name -}}
(
{{- range $i,$arg := $element.Arguments -}}
	{{if $i}}, {{end}}{{ $arg -}}
{{ end -}}
{{ if $element.IsVariadic }}...{{ end -}}
)
{{- $element.Returns -}}
{{ end }}
`

func main() {
	flag.Parse()

	interp := interpreter.MustNewInterpreter()

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

	switch strings.ToLower(*flagOutput) {
	case "json":
		f, err := os.Create("functions.json")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "    ")
		encoder.SetEscapeHTML(false)
		encoder.Encode(fns)
	case "html":
		f, err := os.Create("functions.html")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		t, err := template.New("html").Parse(htmlTemplate)
		if err != nil {
			panic(err)
		}
		t.Execute(f, fns)
	case "md":
		f, err := os.Create("functions.md")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		t, err := texttemplate.New("html").Parse(markdownTemplate)
		if err != nil {
			panic(err)
		}
		t.Execute(f, fns)
	default:
		panic("Unknown output format")
	}
}
