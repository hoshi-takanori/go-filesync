// +build !client

package main

import (
	"net/http"
)

func main() {
	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", SyncHandler)
	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		panic(err)
	}
}

func SyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}

	var msg Message
	err := msg.Decode(r.Body)
	if err != nil || msg.Token != config.Token ||
		(msg.Mode != SyncModeBegin && msg.Mode != SyncModeBoth) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	res := NewMessage(msg.Mode + 1)
	msg.SyncEntries(&res, config.ServerDir)

	err = res.Encode(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
