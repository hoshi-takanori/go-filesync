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

	localDir := CreateLocalDir(now, old)
	remoteDir := CreateRemoteDir(now, old)

	fs, err := remoteDir.List()
	if err != nil {
		panic(err)
	}
	CheckFirstFis(t, now, old, fs)

	fs, err = SyncFiles(SyncModeSend, &localDir, fs)
	if err != nil {
		panic(err)
	}
	CheckSecondFis(t, now, old, fs)

	CheckInterimDir(t, now, old, localDir)

	fs, err = SyncFiles(SyncModeBoth, &remoteDir, fs)
	if err != nil {
		panic(err)
	}
	CheckThirdFis(t, now, old, fs)

	fs, err = SyncFiles(SyncModeReceive, &localDir, fs)
	if err != nil {
		panic(err)
	}
	CheckFis(t, "final", fs)

	CheckFinalDir(t, now, localDir)
	CheckFinalDir(t, now, remoteDir)
}

func CreateLocalDir(now, old time.Time) FakeDir {
	return NewFakeDir("local",
		FakeFInfo("a", now, []byte("aaa")),

		FakeFInfo("b1", now, []byte("b1new")),
		FakeFInfo("b2", now, []byte("b2new")),
		FakeFInfo("b3", now, []byte{}),
		FakeFInfo("b4", now, []byte{}),

		FakeFInfo("c1", old, []byte("c1old")),
		FakeFInfo("c3", old, []byte("c3old")),
	)
}

func CreateRemoteDir(now, old time.Time) FakeDir {
	return NewFakeDir("remote",
		FakeFInfo("a", now, []byte("aaa")),

		FakeFInfo("b1", old, []byte("b1old")),
		FakeFInfo("b3", old, []byte("b3old")),

		FakeFInfo("c1", now, []byte("c1new")),
		FakeFInfo("c2", now, []byte("c2new")),
		FakeFInfo("c3", now, []byte{}),
		FakeFInfo("c4", now, []byte{}),
	)
}

func CheckFirstFis(t *testing.T, now, old time.Time, fis []FInfo) {
	CheckFis(t, "first", fis,
		FakeFInfo("a", now, []byte("aaa")).Copy(),

		FakeFInfo("b1", old, []byte("b1old")).Copy(),
		FakeFInfo("b3", old, []byte("b3old")).Copy(),

		FakeFInfo("c1", now, []byte("c1new")).Copy(),
		FakeFInfo("c2", now, []byte("c2new")).Copy(),
		FakeFInfo("c3", now, []byte{}).Copy(),
		FakeFInfo("c4", now, []byte{}).Copy(),
	)
}

func CheckSecondFis(t *testing.T, now, old time.Time, fis []FInfo) {
	CheckFis(t, "second", fis,
		FakeFInfo("a", now, []byte("aaa")).Copy(),

		FakeFInfo("b1", now, []byte("b1new")),
		FakeFInfo("b2", now, []byte("b2new")),
		FakeFInfo("b3", now, []byte{}).Copy(),

		FakeFInfo("c1", old, []byte("c1old")).Copy(),
		FakeFInfo("c3", old, []byte("c3old")).Copy(),
	)
}

func CheckThirdFis(t *testing.T, now, old time.Time, fis []FInfo) {
	CheckFis(t, "third", fis,
		FakeFInfo("c1", now, []byte("c1new")),
		FakeFInfo("c2", now, []byte("c2new")),
		FakeFInfo("c3", now, []byte{}).Copy(),
	)
}

func CheckInterimDir(t *testing.T, now, old time.Time, dir FakeDir) {
	CheckDir(t, dir,
		FakeFInfo("a", now, []byte("aaa")),

		FakeFInfo("b1", now, []byte("b1new")),
		FakeFInfo("b2", now, []byte("b2new")),
		FakeFInfo("b3", now, []byte{}),

		FakeFInfo("c1", old, []byte("c1old")),
		FakeFInfo("c3", old, []byte("c3old")),
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
