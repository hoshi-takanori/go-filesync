// +build sync

package main

import (
	"testing"
)

func TestSyncFiles(t *testing.T) {
	println("TestSyncFiles")

	remoteFs, err := FSDir("dst").List()
	if err != nil {
		panic(err)
	}

	fs, err := SyncFiles(FSDir("src"), remoteFs)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		PrintFInfo(f)
	}
}
