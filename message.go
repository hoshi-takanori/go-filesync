package main

import (
	"encoding/gob"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Message struct {
	Token   string
	Mode    int
	Entries []Entry
}

type Entry struct {
	Name string
	Fs   []FInfo
}

func NewMessage(mode int) Message {
	return Message{
		Token:   config.Token,
		Mode:    mode,
		Entries: []Entry{},
	}
}

func (msg Message) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(msg)
}

func (msg *Message) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(msg)
}

func (msg *Message) AddEntry(name string, fs []FInfo) {
	msg.Entries = append(msg.Entries, Entry{
		Name: name,
		Fs:   fs,
	})
}

func (msg *Message) ExpandEntries(base string) {
	entries := []Entry{}
	for _, entry := range msg.Entries {
		list, err := filepath.Glob(path.Join(base, entry.Name))
		if err == nil {
			for _, name := range list {
				fi, err := os.Stat(name)
				if err == nil && fi.IsDir() &&
					strings.HasPrefix(name, base+"/") &&
					!strings.Contains(name, "/.") {
					name = strings.TrimPrefix(name, base+"/")
					entries = append(entries, Entry{name, nil})
				}
			}
		}
	}
	msg.Entries = entries
}

func (msg Message) SyncEntries(res *Message, base string) {
	for _, entry := range msg.Entries {
		dir := FSDir{path.Join(base, entry.Name), nil}
		fs, err := SyncFiles(msg.Mode, dir, entry.Fs)
		if err == nil && res != nil {
			res.AddEntry(entry.Name, fs)
		}
	}
}
