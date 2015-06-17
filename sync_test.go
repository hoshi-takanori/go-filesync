// +build sync

package main

import (
	"time"

	"testing"
)

func TestSyncFiles(t *testing.T) {
	println("TestSyncFiles")

	now := time.Now()

	localDir := FakeDir{
		[]FInfo{
			FakeFInfo("a", now, []byte("aaa")),
			FakeFInfo("b", now, []byte("bbb")),
		},
	}

	remoteFs, _ := FakeDir{
		[]FInfo{
			FakeFInfo("a", now, []byte("aaa")),
			FakeFInfo("c", now, []byte("ccc")),
		},
	}.List()

	fs, err := SyncFiles(localDir, remoteFs)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		PrintFInfo(f)
	}
}
