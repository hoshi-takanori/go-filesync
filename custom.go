// +build client

package main

import (
	"log"
	"os"
	"os/user"
	"path"
	"strconv"
)

type CustomDir struct {
	FSDir
	baseDir string
	uid     int
	gid     int
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func NewCustomDir(name string, logger *log.Logger) Dir {
	owner, err := user.Lookup(path.Base(name))
	if err != nil {
		logger.Println("lookup failed: " + err.Error())
		return nil
	}
	uid, _ := strconv.Atoi(owner.Uid)
	gid, _ := strconv.Atoi(owner.Gid)

	baseDir := path.Join(owner.HomeDir, "public_html")
	export := path.Join(baseDir, "export")
	if fi, err := os.Stat(export); err == nil && fi.IsDir() {
		return CustomDir{FSDir{export, logger}, baseDir, uid, gid}
	} else {
		return CustomDir{FSDir{baseDir, logger}, baseDir, uid, gid}
	}
}

func (dir CustomDir) Write(f FInfo) error {
	filePath := dir.Path(f.Name)
	basePath := path.Join(dir.baseDir, f.Name)
	flag := filePath == basePath || Exists(filePath) || Exists(basePath)
	err := dir.FSDir.Write(f)
	if err == nil {
		err = os.Chown(filePath, dir.uid, dir.gid)
	}
	if err == nil && !flag {
		err2 := os.Symlink(path.Join("export", f.Name), basePath)
		if err2 == nil {
			os.Lchown(basePath, dir.uid, dir.gid)
		}
	}
	return err
}
