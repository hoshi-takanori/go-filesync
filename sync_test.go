// +build sync

package main

import (
	"time"

	"testing"
)

func TestSyncFiles(t *testing.T) {
	println("TestSyncFiles")

	now := time.Now()
	old := now.Add(-100 * time.Second)

	localDir, remoteDir := CreateDirs(now, old)
	remoteFs, _ := remoteDir.List()

	CheckFirstFis(t, now, old, remoteFs)

	for i, fi := range remoteFs {
		if fi.Name == "c1" || fi.Name == "c2" {
			remoteDir.Read(&remoteFs[i])
		}
	}

	_, err := SyncFiles(&localDir, remoteFs)
	if err != nil {
		panic(err)
	}

	localFs, _ := localDir.List()

	for i, fi := range localFs {
		if fi.Name == "b1" || fi.Name == "b2" {
			localDir.Read(&localFs[i])
		}
	}

	_, err = SyncFiles(&remoteDir, localFs)
	if err != nil {
		panic(err)
	}

	CheckFinalDir(t, now, localDir)
	CheckFinalDir(t, now, remoteDir)
}

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

func CheckFirstFis(t *testing.T, now, old time.Time, fis []FInfo) {
	CheckFis(t, fis,
		FakeFInfo("a", now, []byte("aaa")).Copy(),

		FakeFInfo("b1", old, []byte("b1old")).Copy(),
		FakeFInfo("b3", old, []byte("b3old")).Copy(),

		FakeFInfo("c1", now, []byte("c1new")).Copy(),
		FakeFInfo("c2", now, []byte("c2new")).Copy(),
		FakeFInfo("c3", now, []byte{}).Copy(),
		FakeFInfo("c4", now, []byte{}).Copy(),
	)
}

func CheckFinalDir(t *testing.T, now time.Time, dir FakeDir) {
	CheckDir(t, dir,
		FakeFInfo("a", now, []byte("aaa")),

		FakeFInfo("b1", now, []byte("b1new")),
		FakeFInfo("b2", now, []byte("b2new")),

		FakeFInfo("c1", now, []byte("c1new")),
		FakeFInfo("c2", now, []byte("c2new")),
	)
}
