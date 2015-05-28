// +build client

package main

import (
	"encoding/gob"
	"net/http"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	println("TestClient")

	client := http.Client{Timeout: time.Duration(10 * time.Second)}

	req, err := http.NewRequest("PUT", "http://localhost:8080/", nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var fs []FInfo
	err = gob.NewDecoder(resp.Body).Decode(&fs)
	if err != nil {
		panic(err)
	}

	for _, f := range fs {
		PrintFInfo(f)
	}
}
