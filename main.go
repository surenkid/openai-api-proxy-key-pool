package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"
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
		errorMessage := "Authorization header is missing"
		log.Printf("[Error] %s", errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	if len(authorization) < 9 || authorization[:7] != "Bearer " {
		errorMessage := "Invalid Authorization header format"
		log.Printf("[Error] %s", errorMessage)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	token := authorization[7:]
	log.Printf("Parsed token: %s", token)

	if token[:3] == "ai-" {
		keys, ok := config.Keys[token]
		if !ok {
			errorMessage := `{"error":{"message":"Invalid Token","code":403}}`
			log.Printf("[Error] %s", errorMessage)
			http.Error(w, errorMessage, http.StatusForbidden)
			return
		}

		index, _ := keyIndex.LoadOrStore(token, 0)
		r.Header.Set("Authorization", "Bearer "+keys[index.(int)])

		nextIndex := (index.(int) + 1) % len(keys)
		keyIndex.Store(token, nextIndex)
		log.Printf("Used key: %s, Updated index: %d", keys[index.(int)], nextIndex)
	}

	target, _ := url.Parse("https://api.openai.com")
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.FlushInterval = time.Millisecond * 200
	proxy.ModifyResponse = func(resp *http.Response) error {
		contentType := resp.Header.Get("Content-Type")
		resp.Header.Set("Content-Type", contentType)
		return nil
	}

	proxy.Director = func(req *http.Request) {
		proxy.Director(req)
		req.Header = r.Header
		req.Header.Set("Transfer-Encoding", r.Header.Get("Transfer-Encoding"))
		req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	}

	proxy.ServeHTTP(w, r)
}

func main() {
	loadConfig()

	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":8124", nil))
}
