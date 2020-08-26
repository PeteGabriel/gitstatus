package main

import (
	"testing"

	is2 "github.com/matryer/is"
)

func TestRecursiveScanFolder(t *testing.T) {
	is := is2.New(t)
	folder := "C:\\Users\\pitadealmeidap\\Documents\\projects"
	repos := scanFolder(make([]string, 0), folder)
	is.True(len(repos) > 1)
}


