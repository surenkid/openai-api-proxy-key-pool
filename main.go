package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

type Config struct {
	Keys map[string][]string `json:"keys"`
}

var config Config
var keyIndex sync.Map

func loadConfig() {
	file, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
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
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	loadConfig()

	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":8124", nil))
}
