// +build sync

package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"testing"
)

func TestSyncFiles(t *testing.T) {
	println("TestSyncFiles")

	exec.Command("sh", "-c", strings.Join([]string{
		"rm -rf src dst",
		"mkdir src dst",
		"echo aaa > src/a",
		"cp -p src/a dst",
		"echo bbb > src/b",
		"echo ccc > dst/c",
		"echo eee > src/e",
		"touch -A -0100 src/e",
		"touch dst/e",
		"echo fff > dst/f",
		"touch -A -0200 dst/f",
		"touch src/f",
	}, ";")).Run()

	logger := log.New(os.Stdout, "", log.LstdFlags)
	src := FSDir{"src", logger}
	dst := FSDir{"dst", logger}

	fs, err := SyncFiles(SyncModeBegin, dst, nil)
	if err != nil {
		panic(err)
	}

	fs, err = SyncFiles(SyncModeSend, src, fs)
	if err != nil {
		panic(err)
	}

	fs, err = SyncFiles(SyncModeBoth, dst, fs)
	if err != nil {
		panic(err)
	}

	fs, err = SyncFiles(SyncModeReceive, src, fs)
	if err != nil {
		panic(err)
	}
}
