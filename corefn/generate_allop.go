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

func mustPrint(n int, err error) {
	if err != nil {
		panic(err)
	}
}

func createDirtyFile(buf *bytes.Buffer, astFile *ast.File) {
	mustPrint(fmt.Fprintf(buf, "package %s\n", *pkg))
	mustPrint(io.WriteString(buf, "import \"github.com/talon-one/talang/interpreter\"\n"))

	mustPrint(io.WriteString(buf, "func AllOperations() []interpreter.TaFunction {\nreturn []interpreter.TaFunction{\n"))
	for _, d := range astFile.Decls {
		if g, ok := d.(*ast.GenDecl); ok {
			if strings.EqualFold(g.Tok.String(), "var") && len(g.Specs) > 0 {
				if v, ok := g.Specs[0].(*ast.ValueSpec); ok {
					if len(v.Names) > 0 {
						mustPrint(fmt.Fprintf(buf, "%s, \n", v.Names[0].Name))
					}
				}
			}
		}
	}
	mustPrint(io.WriteString(buf, "}\n"))
	mustPrint(io.WriteString(buf, "}"))
}

func main() {
	flag.Parse()
	if len(*pkg) == 0 {
		panic("pkg not defined")
	}
	f, err := os.Create(fmt.Sprintf("%s_allop.go", *pkg))
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
	createDirtyFile(&buf, astFile)

	// pretty print
	set = token.NewFileSet()
	astFile, err = parser.ParseFile(set, "", buf.String(), 0)
	if err != nil {
		mustPrint(fmt.Println("Failed to parse package:", err))
		os.Exit(1)
	}
	if err := printer.Fprint(f, set, astFile); err != nil {
		panic(err)
	}
}
