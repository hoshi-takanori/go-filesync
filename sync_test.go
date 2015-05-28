// +build sync

package main

import (
	"testing"
)

func TestSyncFiles(t *testing.T) {
	println("TestSyncFiles")

	remoteFs, err := ListFInfo("dst")
	if err != nil {
		panic(err)
	}

	fs, err := SyncFiles("src", remoteFs)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		PrintFInfo(f)
	}
}
