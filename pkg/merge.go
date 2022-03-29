package merge

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	gotoken "go/token"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"istio.io/pkg/log"

	"golang.org/x/tools/go/ast/astutil"
)

type Merge struct {
	BaseVersionFilename    string
	CurrentVersionFilename string
	OtherVersionFilename   string

	MergeFileFlags []string

	NewBase    string
	NewCurrent string
	NewOther   string
}

func NewMerge(current string, base string, other string) *Merge {

	return &Merge{
		BaseVersionFilename:    base,
		CurrentVersionFilename: current,
		OtherVersionFilename:   other,
		NewBase:                base,
		NewCurrent:             current,
		NewOther:               other,
	}
}

func (m *Merge) MergeFile() {
	fset := gotoken.NewFileSet()
	base, err := parser.ParseFile(fset, m.BaseVersionFilename, nil, 'r')
	if err != nil {
		log.Fatal(err)
	}
	fset = gotoken.NewFileSet()
	current, err := parser.ParseFile(fset, m.CurrentVersionFilename, nil, 'r')
	if err != nil {
		log.Fatal(err)
	}
	fset = gotoken.NewFileSet()
	other, err := parser.ParseFile(fset, m.OtherVersionFilename, nil, 'r')
	if err != nil {
		log.Fatal(err)
	}

	unionImports := make([]*ast.ImportSpec, 0)

	for _, file := range []*ast.File{base, current, other} {
		for _, i := range file.Imports {
			unionImports = append(unionImports, i)
		}
	}

	for _, i := range unionImports {
		for _, file := range []*ast.File{base, current, other} {
			if i.Name != nil {
				astutil.AddNamedImport(fset, file, i.Name.Name, strings.Trim(i.Path.Value, "\""))
			} else {
				astutil.AddImport(fset, file, strings.Trim(i.Path.Value, "\""))
			}
			ast.SortImports(fset, file)
		}
	}

	af, err := os.Create(m.NewBase)
	defer af.Close()
	if err != nil {
		log.Fatal(err)
	}
	fset = gotoken.NewFileSet()
	if err := format.Node(af, fset, base); err != nil {
		log.Fatal(err)
	}

	cf, err := os.Create(m.NewCurrent)
	defer cf.Close()
	if err != nil {
		log.Fatal(err)
	}
	fset = gotoken.NewFileSet()
	if err := format.Node(cf, fset, current); err != nil {
		log.Fatal(err)
	}

	of, err := os.Create(m.NewOther)
	defer of.Close()
	if err != nil {
		log.Fatal(err)
	}
	fset = gotoken.NewFileSet()
	if err := format.Node(of, fset, other); err != nil {
		log.Fatal(err)
	}

}

func (m *Merge) CallNextBinary() {
	log.Infof("Handing exec off to 'git merge-file'...")
	binary, lookErr := exec.LookPath("git")
	if lookErr != nil {
		panic(lookErr)
	}
	if err := syscall.Exec(
		binary,
		[]string{
			"git",
			"merge-file",
			m.NewCurrent,
			m.NewBase,
			m.NewOther,
		},
		[]string{}); err != nil {
		log.Errorf(err)
	}
}

func test() {

	fmt.Println("Hello, world")
}
