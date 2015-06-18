// +build !client

package main

import (
	"encoding/gob"
	"net/http"
)

func main() {
	http.HandleFunc("/", FInfoHandler)
	http.ListenAndServe(":8080", nil)
}

func FInfoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "Page Not Found", http.StatusNotFound)
		return
	}

	fs, err := FSDir(".").List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = gob.NewEncoder(w).Encode(fs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
