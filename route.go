package main

import "net/http"

const (
	rootDir = "."
)

func Route(cfg Config) http.Handler {
	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(rootDir)))
	fsHandler = cfg.MwIncrementHits(fsHandler)
	mux.Handle("/app/*", fsHandler)
	mux.HandleFunc("GET /api/healthz", ApiCheckHealth)
	mux.HandleFunc("GET /api/reset", cfg.ApiResetHits)
	mux.HandleFunc("GET /admin/metrics", cfg.AdminGetHitCount)
	mux.HandleFunc("POST /api/chirps", cfg.ApiPostChirp)
	mux.HandleFunc("GET /api/chirps", cfg.ApiGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.ApiGetChirpByID)
	corsMux := MwAddCors(mux)
	return corsMux
}
