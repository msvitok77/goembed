# goembed
goembed is Golang go:embed parse package

[![Go1.16](https://github.com/msvitok77/goembed/workflows/Go1.16/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go116.yml)
[![Go1.17](https://github.com/msvitok77/goembed/workflows/Go1.17/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go117.yml)
[![Go1.18](https://github.com/msvitok77/goembed/workflows/Go1.18/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go118.yml)
[![Go1.19](https://github.com/msvitok77/goembed/workflows/Go1.19/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go119.yml)
[![Go1.20](https://github.com/msvitok77/goembed/workflows/Go1.20/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go120.yml)
[![Go1.21](https://github.com/msvitok77/goembed/workflows/Go1.21/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go121.yml)
[![Go1.22](https://github.com/msvitok77/goembed/workflows/Go1.22/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go122.yml)
[![Go1.23](https://github.com/msvitok77/goembed/workflows/Go1.23/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go123.yml)
[![Go1.24](https://github.com/msvitok77/goembed/workflows/Go1.24/badge.svg)](https://github.com/msvitok77/goembed/actions/workflows/go124.yml)


### demo
```
package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/msvitok77/goembed"
)

func main() {
	pkg, err := build.Import("github.com/msvitok77/goembed", "", 0)
	if err != nil {
		panic(err)
	}
	fset := token.NewFileSet()
	var files []*ast.File
	for _, file := range pkg.TestGoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, 0)
		if err != nil {
			panic(err)
		}
		files = append(files, f)
	}
	ems,err := goembed.CheckEmbed(pkg.TestEmbedPatternPos, fset, files)
	if err != nil {
		panic(err)
	}
	r := goembed.NewResolve()
	for _, em := range ems {
		files, err := r.Load(pkg.Dir, fset, em)
		if err != nil {
			panic(err)
		}
		for _, f := range files {
			fmt.Println(f.Name, f.Data, f.Hash)
		}
	}
}
```
