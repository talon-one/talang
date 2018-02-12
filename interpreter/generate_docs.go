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
	"github.com/talon-one/talang/interpreter/shared"
)

var flagOutput = flag.String("format", "md", "format to use for output")

const htmlTemplate string = `
<html>
	<style type="text/css">
	details {
		font-family: monospace;
		display: block;
		margin: 1rem 0;
		padding: 0;
	}
	p {
		font-family: monospace;
		display: block;
		margin: 0 0 0 .4rem;
		padding: 0;
	}
	</style>
	<h1>Embedded Functions</h1>
	{{ range $index, $element := . }}
		<details>
			<summary><b>{{ $element.Name -}}</b>
			(
			{{- range $i,$arg := $element.Arguments -}}
				{{if $i}}, {{end}}{{ $arg -}}
			{{ end -}}
			{{ if $element.IsVariadic }}...{{ end -}}
			)
			{{- $element.Returns -}}
			</summary>
			<p>{{ $element.Description }}</p>
		</details>
	{{ end }}
</html>
`

const markdownTemplate string = `# Embedded Functions
{{ range $index, $element := . }}
### {{ $element.Name -}}
(
{{- range $i,$arg := $element.Arguments -}}
	{{if $i}}, {{end}}{{ $arg -}}
{{ end -}}
{{ if $element.IsVariadic }}...{{ end -}}
)
{{- $element.Returns }}
    {{ $element.Description }}
{{ end }}
`

func main() {
	flag.Parse()

	interp := interpreter.MustNewInterpreter()

	type fn struct {
		Arguments []string
		Returns   string
		shared.TaSignature
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
			Arguments:   arguments,
			Returns:     returns,
			TaSignature: f,
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
