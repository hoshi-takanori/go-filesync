package main

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Token string

	Address   string
	ServerDir string

	ServerURL string
	ClientDir string
}

var config Config

func LoadConfig(filename string) error {
	str, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(str, &config)
}
