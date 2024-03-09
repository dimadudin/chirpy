package main

type Config struct {
	fsHits int
}

func NewApiConfig() Config {
	return Config{fsHits: 0}
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
