package main

import (
	"log"
	"net/http"
)

const (
	port = ":8080"
)

func main() {
	cfg := NewApiConfig()
	router := Route(cfg)
	server := http.Server{Addr: port, Handler: router}
	log.Fatal(server.ListenAndServe())
}
