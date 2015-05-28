package main

import (
	"encoding/gob"
	"net/http"
	"os"
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

	fs, err := ListFInfo(".", func(os.FileInfo) bool { return true })
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
