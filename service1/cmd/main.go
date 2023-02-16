package main

import (
	"log"
	"net/http"

	"service1/internal"
)

const (
	port = ":8082"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/generate-salt", internal.GenerateSalt)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatal(err)
	}
}
