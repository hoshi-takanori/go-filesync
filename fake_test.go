// +build fake fake_sync

package main

import (
	"os"
	"path"
	"reflect"
	"time"

	"testing"
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

func (fi FInfo) Copy() FInfo {
	return FInfo{
		Name:    fi.Name,
		Size:    fi.Size,
		Mode:    fi.Mode,
		ModTime: fi.ModTime,
	}
}

func (dir FakeDir) Path(name string) string {
	return path.Join(dir.name, name)
}

func (dir FakeDir) Log(str string) {
	println(str)
}

func (dir FakeDir) List() ([]FInfo, error) {
	fs := []FInfo{}
	for _, fi := range dir.fmap {
		fs = append(fs, fi.Copy())
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

func CheckDir(t *testing.T, dir FakeDir, fis ...FInfo) {
	if len(dir.fmap) != len(fis) {
		t.Errorf("%s: len %d != %d", dir.name, len(dir.fmap), len(fis))
	}

	m := map[string]bool{}
	for name, _ := range dir.fmap {
		m[name] = true
	}

	for _, fi := range fis {
		f, ok := dir.fmap[fi.Name]
		if !ok {
			t.Errorf("%s: %s is missing", dir.name, fi.Name)
		} else {
			if f.Size != fi.Size {
				t.Errorf("%s: %s size %d != %d", dir.name, f.Name, f.Size, fi.Size)
			}
			if f.ModTime != fi.ModTime {
				t.Errorf("%s: %s mtime differ", dir.name, f.Name)
			}
			if !reflect.DeepEqual(f.Content, fi.Content) {
				t.Errorf("%s: %s content differ", dir.name, f.Name)
			}
			delete(m, fi.Name)
		}
	}

	for name, _ := range m {
		t.Errorf("%s: %s exists", dir.name, name)
	}
}

func CheckFis(t *testing.T, name string, fis1 []FInfo, fis2 ...FInfo) {
	CheckDir(t, NewFakeDir(name, fis1...), fis2...)
}

func TestFakeDir(t *testing.T) {
	println("TestFakeDir")

	now := time.Now()
	old := now.Add(-100 * time.Second)

	dir := NewFakeDir("dir",
		FakeFInfo("a", now, []byte("aaa")),
		FakeFInfo("b", old, []byte("bbbbb")),
		FakeFInfo("c", now, []byte{}),
	)

	CheckDir(t, dir,
		FakeFInfo("a", now, []byte("aaa")),
		FakeFInfo("b", old, []byte("bbbbb")),
		FakeFInfo("c", now, []byte{}),
	)

	fis, _ := dir.List()
	CheckFis(t, "fis", fis,
		FakeFInfo("a", now, []byte("aaa")).Copy(),
		FakeFInfo("b", old, []byte("bbbbb")).Copy(),
		FakeFInfo("c", now, []byte{}).Copy(),
	)
}
