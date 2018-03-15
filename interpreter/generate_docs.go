// +build -ignore

package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"html/template"
	texttemplate "text/template"

	"github.com/talon-one/talang"
	"github.com/talon-one/talang/interpreter"
)

var flagOutput = flag.String("format", "md", "format to use for output")
var flagDir = flag.String("dir", ".", "output dir")

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
			<code>{{- $element.Example -}}</code>
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
` + "```lisp" + `
{{ TrimSpace $element.Example }}
` + "```" + `
{{ end }}
`

func main() {
	flag.Parse()

	interp := talang.MustNewInterpreter()

	type fn struct {
		Arguments []string
		Returns   string
		interpreter.TaFunction
	}

	fns := make([]fn, len(interp.Functions))

	for i := 0; i < len(interp.Functions); i++ {
		f := interp.Functions[i]
		arguments := make([]string, len(f.Arguments))
		for j, a := range f.Arguments {
			arguments[j] = a.String()
		}

		packageName := getPackageName(f.Func)
		// some warnings if some data is missing
		if len(strings.TrimSpace(f.Description)) <= 0 {
			log.Printf("WARNING: func `%s' has no `Description`", filepath.Join(packageName, f.Name))
		}
		if len(strings.TrimSpace(f.Example)) <= 0 {
			log.Printf("WARNING: func `%s' has no `Example`", filepath.Join(packageName, f.Name))
		}
		if strings.IndexRune(f.Description, '\t') >= 0 {
			log.Printf("WARNING: func's `%s' `Description` has a TAB character in it", filepath.Join(packageName, f.Name))
		}
		if strings.IndexRune(f.Example, '\t') >= 0 {
			log.Printf("WARNING: func's `%s' `Example` has a TAB character in it", filepath.Join(packageName, f.Name))
		}

		fns[i] = fn{
			Arguments:  arguments,
			Returns:    f.Returns.String(),
			TaFunction: f,
		}
	}

	sort.Slice(fns, func(i, j int) bool {
		return strings.Compare(fns[i].Name, fns[j].Name) < 0
	})

	switch strings.ToLower(*flagOutput) {
	case "json":
		f, err := os.Create(filepath.Join(*flagDir, "functions.json"))
		if err != nil {
			panic(err)
		}
		defer f.Close()
		encoder := json.NewEncoder(f)
		encoder.SetIndent("", "    ")
		encoder.SetEscapeHTML(false)
		err = encoder.Encode(fns)
		if err != nil {
			panic(err)
		}
	case "html":
		f, err := os.Create(filepath.Join(*flagDir, "functions.html"))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		t, err := template.New("html").Funcs(template.FuncMap{"TrimSpace": strings.TrimSpace}).Parse(htmlTemplate)
		if err != nil {
			panic(err)
		}
		err = t.Execute(f, fns)
		if err != nil {
			panic(err)
		}
	case "md":
		f, err := os.Create(filepath.Join(*flagDir, "functions.md"))
		if err != nil {
			panic(err)
		}
		defer f.Close()

		t, err := texttemplate.New("html").Funcs(texttemplate.FuncMap{"TrimSpace": strings.TrimSpace}).Parse(markdownTemplate)
		if err != nil {
			panic(err)
		}
		err = t.Execute(f, fns)
		if err != nil {
			panic(err)
		}
	default:
		panic("Unknown output format")
	}
}

func getPackageName(f interface{}) string {
	args := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), "/")
	for i := len(args) - 1; i >= 0; i-- {
		if len(args[i]) > 0 {
			pos := strings.Index(args[i], ".")
			if pos > -1 {
				return args[i][:pos]
			}
		}
	}
	return ""
}
