package main

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

type FInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	Content []byte
}

func NewFInfo(dir string, fi os.FileInfo, readFlag bool) FInfo {
	f := FInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    fi.Mode(),
		ModTime: fi.ModTime(),
	}

	if readFlag {
		r, err := os.Open(path.Join(dir, f.Name))
		if err == nil {
			f.Content, err = ioutil.ReadAll(r)
		}
	}

	return f
}

func ListFInfo(dir string, readContent func(os.FileInfo) bool) ([]FInfo, error) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return []FInfo{}, err
	}

	fs := []FInfo{}
	for _, fi := range fis {
		if fi.Mode()&os.ModeSymlink != 0 {
			li, err := os.Stat(path.Join(dir, fi.Name()))
			if err == nil {
				fi = li
			}
		}
		if fi.Mode().IsRegular() && !strings.HasPrefix(fi.Name(), ".") {
			fs = append(fs, NewFInfo(dir, fi, readContent(fi)))
		}
	}
	return fs, nil
}
