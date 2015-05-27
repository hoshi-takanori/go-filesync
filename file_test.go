package main

import (
	"testing"
)

func TestListFInfo(t *testing.T) {
	fs, err := ListFInfo(".")
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		println(f.Mode.String(), f.Name, f.ModTime.String(), f.Size)
	}
}
