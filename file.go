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
}

func ListFInfo(dirname string) ([]FInfo, error) {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		return []FInfo{}, err
	}

	fs := []FInfo{}
	for _, fi := range fis {
		if fi.Mode()&os.ModeSymlink != 0 {
			li, err := os.Stat(path.Join(dirname, fi.Name()))
			if err == nil {
				fi = li
			}
		}
		if fi.Mode().IsRegular() && !strings.HasPrefix(fi.Name(), ".") {
			fs = append(fs, FInfo{
				Name:    fi.Name(),
				Size:    fi.Size(),
				Mode:    fi.Mode(),
				ModTime: fi.ModTime(),
			})
		}
	}
	return fs, nil
}
