// +build -ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
)

var pkg = flag.String("pkg", "", "package name")

func mustPrint(n int, err error) {
	if err != nil {
		panic(err)
	}
}

func createDirtyFile(buf *bytes.Buffer, imports []string) {
	mustPrint(fmt.Fprintf(buf, "package %s\n", *pkg))
	mustPrint(io.WriteString(buf, "import (\n"))
	for i := 0; i < len(imports); i++ {
		mustPrint(fmt.Fprintf(buf, "_ \"github.com/talon-one/talang/corefn/%s\"\n", imports[i]))
	}
	mustPrint(io.WriteString(buf, ")\n"))
}

func main() {
	flag.Parse()
	if len(*pkg) == 0 {
		panic("pkg not defined")
	}
	f, err := os.Create(fmt.Sprintf("%s_corefn.go", *pkg))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fileInfo, err := ioutil.ReadDir("./corefn")
	if err != nil {
		panic(err)
	}

	var dirs []string

	for i := 0; i < len(fileInfo); i++ {
		if fileInfo[i].IsDir() {
			dirs = append(dirs, fileInfo[i].Name())
		}
	}

	var buf bytes.Buffer
	createDirtyFile(&buf, dirs)

	// pretty print
	set := token.NewFileSet()
	astFile, err := parser.ParseFile(set, "", buf.String(), 0)
	if err != nil {
		mustPrint(fmt.Println("Failed to parse package:", err))
		os.Exit(1)
	}
	if err := printer.Fprint(f, set, astFile); err != nil {
		panic(err)
	}
}
