# esub

[![CI](https://github.com/winebarrel/esub/actions/workflows/ci.yml/badge.svg)](https://github.com/winebarrel/esub/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/winebarrel/esub.svg)](https://pkg.go.dev/github.com/winebarrel/esub)
[![Go Report Card](https://goreportcard.com/badge/github.com/winebarrel/esub)](https://goreportcard.com/report/github.com/winebarrel/esub)

A library to replace string placeholders with environment variables.

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/winebarrel/esub"
)

func main() {
	os.Setenv("foo", "ZOO")
	os.Setenv("BAR", "baz")
	tmpl := "foo:${foo} BAR:${BAR}"
	out, err := esub.Fill(tmpl)
	if err != nil {
		panic(err)
	}
	fmt.Println(out) //=> "foo:ZOO BAR:baz"
}
```

### Escape Placeholders

```go
tmpl := "foo:$${foo} BAR:${BAR}"
//           ^^^
out, _ := esub.Fill(tmpl)
fmt.Println(out) //=> "foo:${foo} BAR:baz"
```
