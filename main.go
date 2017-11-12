package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math"
	"os"
)

type Reports []Report

type Report struct {
	name string
	line int
	abc  ABC
}

type ABC struct {
	A int
	B int
	C int
}

func main() {
	var path string

	flag.StringVar(&path, "path", "", "Path to file")

	flag.Parse()

	if path == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	reports := reportFile(path)
	fmt.Print(reports)
}

func reportFile(path string) Reports {
	fset := token.NewFileSet()
	n, err := parser.ParseFile(fset, path, nil, 0)

	if err != nil {
		log.Fatal(err)
	}

	var reports Reports

	ast.Inspect(n, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			report := Report{name: fn.Name.Name, line: fset.Position(fn.Pos()).Line}

			ast.Inspect(n, func(n ast.Node) bool {
				reportNode(&report, n)
				return true
			})

			reports = append(reports, report)
			return false
		}
		return true
	})

	return reports
}

func reportNode(report *Report, n ast.Node) {
	switch n := n.(type) {
	case *ast.AssignStmt, *ast.IncDecStmt:
		report.abc.A += 1
	case *ast.CallExpr:
		report.abc.B += 1
	case *ast.IfStmt:
		if n.Else != nil {
			report.abc.C += 1
		}
	case *ast.BinaryExpr, *ast.CaseClause:
		report.abc.C += 1
	}
}

func (reports Reports) String() string {
	var buffer bytes.Buffer

	for _, report := range reports {
		buffer.WriteString(report.String())
		buffer.WriteString("\n")
	}

	return buffer.String()
}

func (report Report) String() string {
	return fmt.Sprintf(
		"%s %d %d",
		report.name,
		report.line,
		report.abc.calc(),
	)
}

func (abc ABC) calc() int {
	a := math.Pow(float64(abc.A), 2)
	b := math.Pow(float64(abc.B), 2)
	c := math.Pow(float64(abc.C), 2)

	return int(math.Sqrt(a + b + c))
}
