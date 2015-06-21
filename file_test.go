// +build file

package main

import (
	"testing"
)

func TestList(t *testing.T) {
	println("TestList")

	dir := FSDir{".", nil}
	fs, err := dir.List()
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		PrintFInfo(f)
	}
}

func PrintFInfo(f FInfo) {
	time := f.ModTime.Format("2006-01-02 15:04:05")
	content := ""
	if f.Content != nil {
		content += "\""
		for i := 0; i < 16 && i < len(f.Content); i++ {
			b := f.Content[i]
			if b >= 0x20 && b < 0x7f {
				content += string([]byte{b})
			} else {
				content += "."
			}
		}
		content += "\""
	}
	println(f.Mode.String(), time, f.Name, f.Size, content)
}
