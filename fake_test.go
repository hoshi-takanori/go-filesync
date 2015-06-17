// +build sync

package main

import (
	"os"
	"time"
)

type FakeDir struct {
	name string
	fmap map[string]FInfo
}

func NewFakeDir(name string, fis ...FInfo) FakeDir {
	dir := FakeDir{name, map[string]FInfo{}}
	for _, fi := range fis {
		dir.fmap[fi.Name] = fi
	}
	return dir
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
	for _, fi := range dir.fmap {
		fs = append(fs, CopyFInfo(fi))
	}
	return fs, nil
}

func (dir FakeDir) Read(f *FInfo) error {
	fi, ok := dir.fmap[f.Name]
	if ok {
		f.Content = fi.Content
		return nil
	}
	return os.ErrNotExist
}

func (dir *FakeDir) Write(f FInfo) error {
	if f.Size == 0 || len(f.Content) == 0 {
		return os.ErrInvalid
	}
	dir.fmap[f.Name] = f
	return nil
}

func (dir *FakeDir) Remove(name string) error {
	_, ok := dir.fmap[name]
	if ok {
		delete(dir.fmap, name)
		return nil
	}
	return os.ErrNotExist
}
