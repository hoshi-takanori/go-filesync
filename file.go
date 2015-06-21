package main

import (
	"io/ioutil"
	"log"
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
	Path(name string) string
	Log(str string)
	List() ([]FInfo, error)
	Read(f *FInfo) error
	Write(f FInfo) error
	Remove(name string) error
}

type FSDir struct {
	name   string
	logger *log.Logger
}

func NewFInfo(fi os.FileInfo) FInfo {
	return FInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    fi.Mode(),
		ModTime: fi.ModTime(),
	}
}

func ListFiles(base string, fis []os.FileInfo, callback func(os.FileInfo)) {
	for _, fi := range fis {
		if fi.Mode()&os.ModeSymlink != 0 {
			li, err := os.Stat(path.Join(base, fi.Name()))
			if err == nil {
				fi = li
			}
		}
		if !strings.HasPrefix(fi.Name(), ".") {
			callback(fi)
		}
	}
}

func (dir FSDir) Path(name string) string {
	return path.Join(dir.name, name)
}

func (dir FSDir) Log(str string) {
	if dir.logger != nil {
		dir.logger.Println(str)
	}
}

func (dir FSDir) List() ([]FInfo, error) {
	fis, err := ioutil.ReadDir(dir.name)
	if err != nil {
		return nil, err
	}

	fs := []FInfo{}
	ListFiles(dir.name, fis, func(fi os.FileInfo) {
		if fi.Mode().IsRegular() {
			fs = append(fs, NewFInfo(fi))
		}
	})
	return fs, nil
}

func (dir FSDir) Read(f *FInfo) error {
	r, err := os.Open(dir.Path(f.Name))
	if err != nil {
		return err
	}
	defer r.Close()

	f.Content, err = ioutil.ReadAll(r)
	return err
}

func (dir FSDir) Write(f FInfo) error {
	err := ioutil.WriteFile(dir.Path(f.Name), f.Content, 0644)
	if err == nil {
		err = os.Chtimes(dir.Path(f.Name), f.ModTime, f.ModTime)
	}
	return err
}

func (dir FSDir) Remove(name string) error {
	return os.Remove(dir.Path(name))
}
