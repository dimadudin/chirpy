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
	mux.HandleFunc("POST /api/refresh", cfg.ApiRefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.ApiRevokeToken)
	mux.HandleFunc("GET /admin/metrics", cfg.AdminGetHitCount)

	mux.HandleFunc("POST /api/users", cfg.ApiCreateUser)
	mux.HandleFunc("PUT /api/users", cfg.ApiUpdateUser)
	mux.HandleFunc("POST /api/login", cfg.ApiLogin)

	mux.HandleFunc("POST /api/chirps", cfg.ApiPostChirp)
	mux.HandleFunc("GET /api/chirps", cfg.ApiGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.ApiGetChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.ApiDeleteChirpByID)

	return MwAddCors(mux)
}
