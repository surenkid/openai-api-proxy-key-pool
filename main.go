package main

import (
	"log"
	"net/http"
)

func main() {
	config, err := LoadConfig("config/config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	http.HandleFunc("/", ProxyHandler(config))
	log.Fatal(http.ListenAndServe(":8124", nil))
}
