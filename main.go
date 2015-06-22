// +build !client

package main

import (
	"log"
	"net/http"
	"os"
	"path"
)

var logger *log.Logger

func main() {
	err := LoadConfig("config.json")
	if err != nil {
		panic(err)
	}

	logger = log.New(os.Stdout, "", log.LstdFlags)

	http.HandleFunc("/", SyncHandler)
	err = http.ListenAndServe(config.Address, nil)
	if err != nil {
		panic(err)
	}
}

func SyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		logger.Println("bad method: " + r.Method)
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}

	var msg Message
	err := msg.Decode(r.Body)
	if err != nil || msg.Token != config.Token ||
		(msg.Mode != SyncModeBegin && msg.Mode != SyncModeBoth) {
		if err != nil {
			logger.Println("decode failed: " + err.Error())
		} else if msg.Token != config.Token {
			logger.Println("bad token: " + msg.Token)
		} else {
			logger.Printf("bad mode: %d\n", msg.Mode)
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	remote := r.Header.Get("X-Forwarded-For")
	if remote == "" {
		remote = r.RemoteAddr
	}
	logger.Printf("mode %d from %s\n", msg.Mode, remote)

	if msg.Mode == SyncModeBegin {
		msg.ExpandEntries(config.ServerDir, logger)
	}

	res := NewMessage(msg.Mode + 1)
	msg.SyncEntries(&res, func(name string) Dir {
		return FSDir{path.Join(config.ServerDir, name), logger}
	})

	err = res.Encode(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
