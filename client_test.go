// +build client

package main

import (
	"net/http"
	"time"

	"testing"
)

func TestClient(t *testing.T) {
	println("TestClient")

	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	client := http.Client{Timeout: time.Duration(10 * time.Second)}

	msg := NewMessage(SyncModeBegin)
	msg.AddEntry("test", nil)

	res, err := SyncClientOne(&client, msg)
	if err != nil {
		panic(err)
	}

	for _, entry := range res.Entries {
		println(entry.Name)
		for _, f := range entry.Fs {
			time := f.ModTime.Format("2006-01-02 15:04:05")
			println(f.Mode.String(), time, f.Name, f.Size)
		}
	}
}
