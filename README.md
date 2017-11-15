# ABCGo

[![Go Report Card](https://goreportcard.com/badge/github.com/droptheplot/abcgo)](https://goreportcard.com/report/github.com/droptheplot/abcgo)
[![GoDoc](https://godoc.org/github.com/droptheplot/abcgo?status.svg)](https://godoc.org/github.com/droptheplot/abcgo)


ABC metrics for Go source code.

## Definition

ABCGo uses these rules to calculate ABC:

* Add one to the **assignment** count when:
  * Occurrence of an assignment operator: `=`, `*=`, `/=`, `%=`, `+=`, `<<=`, `>>=`, `&=`, `^=`.
  * Occurrence of an increment or a decrement operator: `++`, `--`.
* Add one to **branch** count when:
  * Occurrence of a function call.
* Add one to **condition** count when:
  * Occurrence of a conditional operator: `<`, `>`, `<=`, `>=`, `==`, `!=`.
  * Occurrence of the following keywords: `else`, `case`.

Final score is calculated as follows:

<img src="https://wikimedia.org/api/rest_v1/media/math/render/svg/871176d94f9d4a290ba3c479b24b815567e1eaa1" />

[Read more about ABC metrics.](https://en.wikipedia.org/wiki/ABC_Software_Metric)

## Getting Started

### Installation

```shell
$ go get -u github.com/droptheplot/abcgo
$ (cd $GOPATH/src/github.com/droptheplot/abcgo && go install)
```

### Usage

```shell
$ abcgo -path main.go
main.go:28   init   9    1   8    5
main.go:54   main   13   5   13   1

$ abcgo -path ./
main.go:28        init            9    1   8    5
main.go:54        main            13   5   13   1
main_test.go:54   TestSomething   9    0   9    2

$ abcgo -path main.go -format json
[
  {
    "path": "main.go",
    "line": 54,
    "name": "main",
    "assignment": 5,
    "branch": 13,
    "condition": 1,
    "score": 13
  },
  {
    "path": "main.go",
    "line": 54,
    "name": "init",
    "assignment": 1,
    "branch": 8,
    "condition": 5,
    "score": 9
  }
]
```

#### Options

* `-path` - Path to file or directory.
* `-format` - Output format (`table` (default) or `json`).
* `-sort` - Sort functions by score.

### Plugins

* [Vim](https://github.com/droptheplot/abcgo/vim)
