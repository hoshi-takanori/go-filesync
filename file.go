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

func NewFInfo(dir string, fi os.FileInfo) FInfo {
	return FInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    fi.Mode(),
		ModTime: fi.ModTime(),
	}
}

func (f *FInfo) ReadContent(dir string) error {
	r, err := os.Open(path.Join(dir, f.Name))
	if err != nil {
		return err
	}

	f.Content, err = ioutil.ReadAll(r)
	return err
}

func ListFInfo(dir string) ([]FInfo, error) {
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
			fs = append(fs, NewFInfo(dir, fi))
		}
	}
	return fs, nil
}
