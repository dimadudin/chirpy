package main

import "github.com/dimadudin/web-server-go/internal/database"

type Config struct {
	db     *database.DB
	fsHits int
}

func NewApiConfig(db *database.DB) Config {
	return Config{db: db, fsHits: 0}
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
