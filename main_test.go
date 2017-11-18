package main

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportFile(t *testing.T) {
	fileSet := token.NewFileSet()

	testCases := []struct {
		rule   string
		report Report
		src    string
	}{
		{
			"=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar := 1
				}
			`,
		},
		{
			"*=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar *= 1
				}
			`,
		},
		{
			"/=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar /= 1
				}
			`,
		},
		{
			"%=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar %= 1
				}
			`,
		},
		{
			"+=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar += 1
				}
			`,
		},
		{
			"<<=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar <<= 1
				}
			`,
		},
		{
			">>=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar >>= 1
				}
			`,
		},
		{
			"&=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar &= 1
				}
			`,
		},
		{
			"^=",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar ^= 1
				}
			`,
		},
		{
			"++",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar++
				}
			`,
		},
		{
			"--",
			Report{Assignment: 1},
			`
				package main

				func main() {
					bar--
				}
			`,
		},
		{
			"Function call",
			Report{Branch: 1},
			`
				package main

				func main() {
					bar()
				}
			`,
		},
		{
			"Anonymous function call",
			Report{Branch: 1},
			`
				package main

				func main() {
					func() {}()
				}
			`,
		},
		{
			"Function call inside anonymous function",
			Report{Branch: 1},
			`
				package main

				func main() {
					func() {
						bar()
					}
				}
			`,
		},
		{
			"Function call inside anonymous function call",
			Report{Branch: 2},
			`
				package main

				func main() {
					func() {
						bar()
					}()
				}
			`,
		},
		{
			">",
			Report{Condition: 1},
			`
				package main

				func main() {
					if 1 > 2 {}
				}
			`,
		},
		{
			"<",
			Report{Condition: 1},
			`
				package main

				func main() {
					if 1 < 2 {}
				}
			`,
		},
		{
			"<=",
			Report{Condition: 1},
			`
				package main

				func main() {
					if 1 <= 2 {}
				}
			`,
		},
		{
			">=",
			Report{Condition: 1},
			`
				package main

				func main() {
					if 1 >= 2 {}
				}
			`,
		},
		{
			"==",
			Report{Condition: 1},
			`
				package main

				func main() {
					if 1 == 2 {}
				}
			`,
		},
		{
			"!=",
			Report{Condition: 1},
			`
				package main

				func main() {
					if 1 != 2 {}
				}
			`,
		},
		{
			"`if` without `else`",
			Report{},
			`
				package main

				func main() {
					if true {}
				}
			`,
		},
		{
			"`if` with `else`",
			Report{Condition: 1},
			`
				package main

				func main() {
					if true {} else {}
				}
			`,
		},
		{
			"`case`",
			Report{Condition: 2},
			`
				package main

				func main() {
					switch nil {
						case true:
						case false:
					}
				}
			`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.rule, func(t *testing.T) {
			node, _ := parser.ParseFile(fileSet, "", tc.src, 0)
			report := reportFile(fileSet, node)[0]

			assert.Equal(t, tc.report.Assignment, report.Assignment, "Assignment")
			assert.Equal(t, tc.report.Branch, report.Branch, "Branch")
			assert.Equal(t, tc.report.Condition, report.Condition, "Condition")
		})
	}
}
