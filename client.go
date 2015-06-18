// +build client

package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func main() {
	client := http.Client{Timeout: time.Duration(10 * time.Second)}

	fis, err := ioutil.ReadDir("base")
	if err != nil {
		panic(err)
	}

	msg := NewMessage(SyncModeBegin)
	for _, fi := range fis {
		if fi.IsDir() && !strings.HasPrefix(fi.Name(), ".") {
			msg.AddEntry(fi.Name(), nil)
		}
	}

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
	var buf bytes.Buffer
	err := msg.Encode(&buf)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", "http://localhost:8080/", &buf)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var res Message
	err = res.Decode(resp.Body)
	if err != nil {
		return err
	}

	res.SyncEntries(msg2, "base")

	return nil
}
