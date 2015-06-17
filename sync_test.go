// +build sync

package main

import (
	"time"

	"testing"
)

func CreateDirs(now, old time.Time) (FakeDir, FakeDir) {
	localDir := NewFakeDir("local",
		FakeFInfo("a", now, []byte("aaa")),

		FakeFInfo("b1", now, []byte("b1new")),
		FakeFInfo("b2", now, []byte("b2new")),
		FakeFInfo("b3", now, []byte{}),
		FakeFInfo("b4", now, []byte{}),

		FakeFInfo("c1", old, []byte("c1old")),
		FakeFInfo("c3", old, []byte("c3old")),
	)

	remoteDir := NewFakeDir("remote",
		FakeFInfo("a", now, []byte("aaa")),

		FakeFInfo("b1", old, []byte("b1old")),
		FakeFInfo("b3", old, []byte("b3old")),

		FakeFInfo("c1", now, []byte("c1new")),
		FakeFInfo("c2", now, []byte("c2new")),
		FakeFInfo("c3", now, []byte{}),
		FakeFInfo("c4", now, []byte{}),
	)

	return localDir, remoteDir
}

func TestSyncFiles(t *testing.T) {
	println("TestSyncFiles")

	now := time.Now()
	old := now.Add(-100 * time.Second)

	localDir, remoteDir := CreateDirs(now, old)
	remoteFs, _ := remoteDir.List()

	fs, err := SyncFiles(&localDir, remoteFs)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		PrintFInfo(f)
	}
}
