package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/dimadudin/web-server-go/internal/database"
	"github.com/joho/godotenv"
)

const (
	port   = "8080"
	dbPath = "./database.json"
)

func main() {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	if *dbg {
		os.Remove(dbPath)
	}

	db, err := database.NewDB(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaApiKey := os.Getenv("POLKA_API_KEY")

	cfg := NewApiConfig(db, jwtSecret, polkaApiKey)

	router := Route(cfg)
	server := http.Server{Addr: ":" + port, Handler: router}
	log.Fatal(server.ListenAndServe())
}
