// +build sync

package main

import (
	"os"
	"time"
)

type FakeDir struct {
	fis []FInfo
}

func FakeFInfo(name string, mtime time.Time, content []byte) FInfo {
	return FInfo{
		Name:    name,
		Size:    int64(len(content)),
		Mode:    0644,
		ModTime: mtime,
		Content: content,
	}
}

func CopyFInfo(fi FInfo) FInfo {
	return FInfo{
		Name:    fi.Name,
		Size:    fi.Size,
		Mode:    fi.Mode,
		ModTime: fi.ModTime,
	}
}

func (dir FakeDir) List() ([]FInfo, error) {
	fs := []FInfo{}
	for _, fi := range dir.fis {
		fs = append(fs, CopyFInfo(fi))
	}
	return fs, nil
}

func (dir FakeDir) Read(f *FInfo) error {
	for _, fi := range dir.fis {
		f.Content = fi.Content
		return nil
	}
	return os.ErrNotExist
}
