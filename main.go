package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"path/filepath"
	"sort"
	"text/tabwriter"
)

// Reports is a collection of Report.
type Reports []Report

// Report contains statistics for single function.
type Report struct {
	Path       string `json:"path"`
	Line       int    `json:"line"`
	Name       string `json:"name"`
	Assignment int    `json:"assignment"`
	Branch     int    `json:"branch"`
	Condition  int    `json:"condition"`
	Score      int    `json:"score"`
}

func main() {
	var (
		path    string
		format  string
		sort    bool
		noTest  bool
		reports Reports
	)

	flag.StringVar(&path, "path", "", "Path to file")
	flag.StringVar(&format, "format", "table", "Output format")
	flag.BoolVar(&sort, "sort", false, "Sort by score")
	flag.BoolVar(&noTest, "no-test", false, "Skip *_test.go files")

	flag.Parse()

	if path == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	fileList := listFiles(path)
	fileSet := token.NewFileSet()

	for _, path := range fileList {
		node, err := parser.ParseFile(fileSet, path, nil, 0)

		if err != nil {
			continue
		}

		if noTest && isTest(path) {
			continue
		}

		reports = append(reports, reportFile(fileSet, node)...)
	}

	if sort {
		reports.Sort()
	}

	switch format {
	case "summary":
		reports.renderSummary()
	case "table":
		reports.renderTable()
	case "json":
		reports.renderJSON()
	case "raw":
		reports.renderRaw()
	default:
		panic("unknown format.")
	}
}

func listFiles(path string) []string {
	var fileList []string

	fileInfo, _ := os.Stat(path)

	appendAbsPath := func(p string) {
		p, _ = filepath.Abs(p)
		fileList = append(fileList, p)
	}

	if fileInfo.IsDir() {
		filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
			if filepath.Ext(f.Name()) == ".go" {
				appendAbsPath(p)
			}

			return nil
		})
	} else {
		appendAbsPath(fileInfo.Name())
	}

	return fileList
}

func reportFile(fset *token.FileSet, n ast.Node) Reports {
	var reports Reports

	ast.Inspect(n, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			report := Report{
				Path: fset.File(fn.Pos()).Name(),
				Line: fset.Position(fn.Pos()).Line,
				Name: fn.Name.Name,
			}

			ast.Inspect(n, func(n ast.Node) bool {
				reportNode(&report, n)
				return true
			})

			report.Calc()
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
		report.Assignment++
	case *ast.CallExpr:
		report.Branch++
	case *ast.IfStmt:
		if n.Else != nil {
			report.Condition++
		}
	case *ast.BinaryExpr, *ast.CaseClause:
		report.Condition++
	}
}

func (report Report) String() string {
	return fmt.Sprintf(
		"%s:%d\t%s\t%d\t%d\t%d\t%d",
		report.Path,
		report.Line,
		report.Name,
		report.Score,
		report.Assignment,
		report.Branch,
		report.Condition,
	)
}

// Calc updates Score value.
func (report *Report) Calc() {
	a := math.Pow(float64(report.Assignment), 2)
	b := math.Pow(float64(report.Branch), 2)
	c := math.Pow(float64(report.Condition), 2)

	report.Score = int(math.Sqrt(a + b + c))
}

func (reports Reports) renderSummary() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()
	a, b, c := 0, 0, 0
	for _, report := range reports {
		a = a + report.Assignment
		b = b + report.Branch
		c = c + report.Condition
	}

	fmt.Fprintln(w, "\tA\tB\tC")
	fmt.Fprintf(w, "%s\t%d\t%d\t%d\n", "Project summary:", a, b, c)
}

func (reports Reports) renderTable() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	defer w.Flush()

	fmt.Fprintln(w, "Source\tFunc\tScore\tA\tB\tC")

	for _, report := range reports {
		fmt.Fprintln(w, report.String())
	}
}

func (reports Reports) renderJSON() {
	bytes, err := json.Marshal(reports)

	if err != nil {
		fmt.Println(err)
	}

	os.Stdout.Write(bytes)
}

func (reports Reports) renderRaw() {
	for _, report := range reports {
		fmt.Printf(
			"%s\t%d\t%s\t%d\n",
			report.Path,
			report.Line,
			report.Name,
			report.Score,
		)
	}
}

func (reports Reports) Sort() {
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].Score > reports[j].Score
	})
}

func isTest(p string) bool {
	match, err := filepath.Match("*_test.go", filepath.Base(p))

	return match && err == nil
}
