package main

import (
	"encoding/json"
	"log"
	"net/http"
	"path/filepath"
	"os"
	"sync"

	"github.com/r3labs/sse"
)

type Config struct {
	Keys map[string][]string `json:"keys"`
}

var config Config
var keyIndex sync.Map

func loadConfig() {
	configPath := "config/config.json"
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
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	if len(authorization) == 0 {
		http.Error(w, "Authorization header is missing", http.StatusBadRequest)
		return
	}

	if len(authorization) < 9 || authorization[:7] != "Bearer " {
		http.Error(w, "Invalid Authorization header format", http.StatusBadRequest)
		return
	}

	token := authorization[7:]
	if token[:3] == "ai-" {
		keys, ok := config.Keys[token]
		if !ok {
			http.Error(w, `{"error":{"message":"Invalid Token","code":403}}`, http.StatusForbidden)
			return
		}

		index, _ := keyIndex.LoadOrStore(token, 0)
		r.Header.Set("Authorization", "Bearer "+keys[index.(int)])

		nextIndex := (index.(int) + 1) % len(keys)
		keyIndex.Store(token, nextIndex)
	}

	proxyURL := "https://api.openai.com" + r.RequestURI
	req, err := http.NewRequest(r.Method, proxyURL, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	client := sse.NewClient(proxyURL)
	client.Connection.Transport = http.DefaultTransport

	err = client.SubscribeRaw(req, func(msg *sse.Event) {
		if _, err := w.Write(msg.Data); err != nil {
			log.Printf("Error writing event data to response: %v", err)
		}
	})

	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
}

func main() {
	loadConfig()

	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":8124", nil))
}
