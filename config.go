package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Keys     map[string][]string `json:"keys"`
	Helicone string              `json:"helicone"`
}

func LoadConfig(configPath string) (config Config, err error) {
	file, err := os.Open(configPath)
	if err != nil {
		cwd, _ := os.Getwd()
		fullPath := filepath.Join(cwd, configPath)
		log.Fatalf("Error opening config file at path %s: %v", fullPath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Error decoding config file:", err)
	}
	return
}
