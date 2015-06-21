// +build client

package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main() {
	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	client := http.Client{Timeout: time.Duration(10 * time.Second)}

	fis, err := ioutil.ReadDir(config.ClientDir)
	if err != nil {
		panic(err)
	}

	msg := NewMessage(SyncModeBegin)
	ListFiles(config.ClientDir, fis, func(fi os.FileInfo) {
		if fi.IsDir() {
			msg.AddEntry(fi.Name(), nil)
		}
	})

	msg2 := NewMessage(SyncModeBoth)
	err = SyncClient(&client, msg, &msg2)
	if err != nil {
		panic(err)
	}

	err = SyncClient(&client, msg2, nil)
	if err != nil {
		panic(err)
	}
}

func SyncClient(client *http.Client, msg Message, msg2 *Message) error {
	res, err := SyncClientOne(client, msg)
	if err != nil {
		return err
	}

	res.SyncEntries(msg2, config.ClientDir)

	return nil
}

func SyncClientOne(client *http.Client, msg Message) (*Message, error) {
	var buf bytes.Buffer
	err := msg.Encode(&buf)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", config.ServerURL, &buf)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Bad Response: " + resp.Status)
	}

	var res Message
	err = res.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
