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

type Dir interface {
	List() ([]FInfo, error)
	Read(f *FInfo) error
}

type FSDir string

func NewFInfo(fi os.FileInfo) FInfo {
	return FInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    fi.Mode(),
		ModTime: fi.ModTime(),
	}
}

func (dir FSDir) List() ([]FInfo, error) {
	fis, err := ioutil.ReadDir(string(dir))
	if err != nil {
		return []FInfo{}, err
	}

	fs := []FInfo{}
	for _, fi := range fis {
		if fi.Mode()&os.ModeSymlink != 0 {
			li, err := os.Stat(path.Join(string(dir), fi.Name()))
			if err == nil {
				fi = li
			}
		}
		if fi.Mode().IsRegular() && !strings.HasPrefix(fi.Name(), ".") {
			fs = append(fs, NewFInfo(fi))
		}
	}
	return fs, nil
}

func (dir FSDir) Read(f *FInfo) error {
	r, err := os.Open(path.Join(string(dir), f.Name))
	if err != nil {
		return err
	}

	f.Content, err = ioutil.ReadAll(r)
	return err
}
