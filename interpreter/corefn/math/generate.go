// +build -ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"os"
	"strings"
)

var pkg = flag.String("pkg", "", "package name")

func main() {
	flag.Parse()
	if len(*pkg) == 0 {
		panic("pkg not defined")
	}
	f, err := os.Create(fmt.Sprintf("%s_generated.go", *pkg))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	set := token.NewFileSet()
	astFile, err := parser.ParseFile(set, fmt.Sprintf("%s.go", *pkg), nil, 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	fmt.Fprintf(&buf, "package %s\n", *pkg)
	io.WriteString(&buf, "import \"github.com/talon-one/talang/interpreter/shared\"\n")

	io.WriteString(&buf, "func AllOperations() []shared.TaSignature {\nreturn []shared.TaSignature{\n")
	for _, d := range astFile.Decls {
		if g, ok := d.(*ast.GenDecl); ok {
			if strings.EqualFold(g.Tok.String(), "var") && len(g.Specs) > 0 {
				if v, ok := g.Specs[0].(*ast.ValueSpec); ok {
					if len(v.Names) > 0 {
						fmt.Fprintf(&buf, "%s, \n", v.Names[0].Name)
					}
				}
			}
		}
	}
	io.WriteString(&buf, "}\n")
	io.WriteString(&buf, "}")

	// pretty print
	set = token.NewFileSet()
	astFile, err = parser.ParseFile(set, "", buf.String(), 0)
	if err != nil {
		fmt.Println("Failed to parse package:", err)
		os.Exit(1)
	}
	printer.Fprint(f, set, astFile)
}
