//go:build go1.16
// +build go1.16

package goembed

import (
	"embed"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unsafe"
)

//go:embed testdata/data1.txt
var data1 string

var (
	//go:embed testdata/data2.txt
	data2 []byte

	//go:embed testdata
	fs embed.FS
)

type file struct {
	name string
	data string
	hash [16]byte
}

type myfs struct {
	files *[]file
}

func TestEmbed(t *testing.T) {
	if data1 != "hello data1" {
		t.Fail()
	}
	if string(data2) != "hello data2" {
		t.Fail()
	}
	files := *(*myfs)(unsafe.Pointer(&fs)).files
	for _, file := range files {
		t.Log(file.name, file.data, file.hash)
	}
}

func TestResolve(t *testing.T) {
	pkg, err := build.Import("github.com/msvitok77/goembed", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	var files []*ast.File
	for _, file := range pkg.TestGoFiles {
		f, err := parser.ParseFile(fset, filepath.Join(pkg.Dir, file), nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, f)
	}
	ems, err := CheckEmbed(pkg.TestEmbedPatternPos, fset, files)
	if err != nil {
		t.Fatal(err)
	}
	r := NewResolve()
	var checkData1 bool
	var checkData2 bool
	var checkFS bool
	for _, em := range ems {
		files, err := r.Load(pkg.Dir, fset, em)
		if err != nil {
			t.Fatal("error load", em, err)
		}
		switch em.Name {
		case "data1":
			checkData1 = true
			if string(files[0].Data) != "hello data1" {
				t.Fail()
			}
		case "data2":
			checkData2 = true
			if string(files[0].Data) != "hello data2" {
				t.Fail()
			}
		}
		if em.Kind == EmbedFiles && em.Name == "fs" {
			checkFS = true
			files = BuildFS(files)
			for _, f := range files {
				t.Log(f.Name, string(f.Data), f.Hash)
			}
			var info1 []string
			var info2 []string
			mfiles := *(*myfs)(unsafe.Pointer(&fs)).files

			runtimeVersion := runtime.Version()[:6]
			switch runtimeVersion {
			default:
				info1 = mfilesInfos(mfiles, true)
				info2 = filesInfos(files, true)
			case "go1.19", "go1.20", "go1.21", "go1.22", "go1.23":
				t.Log("go1.19, go1.20, go1.21, go1.22, go1.23 compiler uses notsha256.Sum256 skip hash check")

				info1 = mfilesInfos(mfiles, false)
				info2 = filesInfos(files, false)
			case "go1.24":
				t.Log("go1.24 compiler uses sha256.Sum256 + sum[0] ^= 0xff skip hash check")
				info1 = mfilesInfos(mfiles, false)
				info2 = filesInfos(files, false)
			}
			if strings.Join(info1, ";") != strings.Join(info2, ";") {
				t.Fatalf("build fs error:\n%v\n%v", info1, info2)
			}
		}
	}
	if !checkData1 || !checkData2 || !checkFS {
		t.Fatal("not found embed", checkData1, checkData2, checkFS)
	}
}

func mfilesInfos(mfiles []file, withHash bool) []string {
	var info []string
	for _, file := range mfiles {
		if withHash {
			info = append(info, fmt.Sprintf("%v,%v,%v", file.name, file.data, file.hash))
		} else {
			info = append(info, fmt.Sprintf("%v,%v", file.name, file.data))
		}
	}
	return info
}

func filesInfos(files []*File, withHash bool) []string {
	var info []string
	for _, f := range files {
		if withHash {
			info = append(info, fmt.Sprintf("%v,%v,%v", f.Name, string(f.Data), f.Hash))
		} else {
			info = append(info, fmt.Sprintf("%v,%v", f.Name, string(f.Data)))
		}
	}
	return info
}
