package templates

import (
	"go/build"
	"path/filepath"
)

// Name is the idiomatic mannerto refer to a TemplateSet
type Name int

// PackagePath finds the full path for the specified
// golang import path
//
// i.e., PackPath("github.com/elos/ehttp")
// => "~/Nick/workspace/go/src/github.com/elos/ehttp" (on my computer)
func PackagePath(importPath string) string {
	p, err := build.Default.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "."
	}
	return p.Dir
}

// Prepend creates a slice of string from variadic
// arguments with the guarantee that the slice will
// be of size >= 1, with index 0 equal to s
//
// Prepend is useful for constructing templateSets
// i.e.,
//		func Root(v ...string) []string {
//			return Prepend("root.tmpl", v...)
//		}
func Prepend(s string, v ...string) []string {
	l := make([]string, len(v)+1)
	l[0] = s
	for i := range v {
		l[i+1] = v[i]
	}
	return l
}

// JoinDir prepends a full path to a slice of relative paths
// i.e., base = "/root/here/", files = []string{"1.go", "2.go"}
// JoinDir(base, files) => []string{"/root/here/1.go", "/root/here/2.go"}
// Useful for building templateSet paths.
func JoinDir(base string, files []string) []string {
	r := make([]string, len(files))
	for i := range files {
		r[i] = filepath.Join(base, files[i])
	}
	return r
}
