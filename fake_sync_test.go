// +build fake_sync

package main

import (
	"time"

	"testing"
)

func TestFakeSyncFiles(t *testing.T) {
	println("TestFakeSyncFiles")

	now := time.Now()
	old := now.Add(-100 * time.Second)

	clientDir := CreateClientDir(now, old)
	serverDir := CreateServerDir(now, old)

	fs, err := SyncFiles(SyncModeBegin, &serverDir, nil)
	if err != nil {
		panic(err)
	}
	CheckFirstFis(t, now, old, fs)

	println("send")
	fs, err = SyncFiles(SyncModeSend, &clientDir, fs)
	if err != nil {
		panic(err)
	}
	CheckSecondFis(t, now, old, fs)

	println("both")
	fs, err = SyncFiles(SyncModeBoth, &serverDir, fs)
	if err != nil {
		panic(err)
	}
	CheckThirdFis(t, now, old, fs)

	println("receive")
	fs, err = SyncFiles(SyncModeReceive, &clientDir, fs)
	if err != nil {
		panic(err)
	}
	CheckFis(t, "final", fs)

	CheckFinalDir(t, now, clientDir)
	CheckFinalDir(t, now, serverDir)
}

func CreateClientDir(now, old time.Time) FakeDir {
	return NewFakeDir("client",
		FakeFInfo("a", now, []byte("aaa")),

		FakeFInfo("b1", now, []byte("b1new")),
		FakeFInfo("b2", now, []byte("b2new")),
		FakeFInfo("b3", now, []byte{}),
		FakeFInfo("b4", now, []byte{}),

		FakeFInfo("c1", old, []byte("c1old")),
		FakeFInfo("c3", old, []byte("c3old")),
	)
}

func CreateServerDir(now, old time.Time) FakeDir {
	return NewFakeDir("server",
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
		FakeFInfo("b4", now, []byte{}).Copy(),

		FakeFInfo("c3", old, []byte("c3old")).Copy(),
	)
}

func CheckThirdFis(t *testing.T, now, old time.Time, fis []FInfo) {
	CheckFis(t, "third", fis,
		FakeFInfo("b3", old, []byte("b3old")).Copy(),

		FakeFInfo("c1", now, []byte("c1new")),
		FakeFInfo("c2", now, []byte("c2new")),
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
