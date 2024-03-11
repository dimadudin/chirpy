package main

import (
	"log"
	"net/http"

	"github.com/dimadudin/web-server-go/internal/database"
)

const (
	port   = "8080"
	dbPath = "./database.json"
)

func main() {
	db, err := database.NewDB(dbPath)
	cfg := NewApiConfig(db)
	if err != nil {
		log.Fatal(err)
	}
	router := Route(cfg)
	server := http.Server{Addr: ":" + port, Handler: router}
	log.Fatal(server.ListenAndServe())
}
