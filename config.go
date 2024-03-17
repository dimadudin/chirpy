package main

import "github.com/dimadudin/web-server-go/internal/database"

type Config struct {
	db          *database.DB
	jwtSecret   string
	polkaApiKey string
	fsHits      int
}

func NewApiConfig(db *database.DB, jwtSecret string, polkaApiKey string) Config {
	return Config{db: db, jwtSecret: jwtSecret, polkaApiKey: polkaApiKey, fsHits: 0}
}

func (cfg *Config) RegisterHit() {
	cfg.fsHits++
}

func (cfg *Config) GetHitCount() int {
	return cfg.fsHits
}

func (cfg *Config) ResetHits() {
	cfg.fsHits = 0
}
