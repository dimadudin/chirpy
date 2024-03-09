package main

import (
	"fmt"
	"net/http"
)

func middlewareCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

type apiConfig struct {
	fsHits int
}

func newApiConfig() apiConfig {
	return apiConfig{fsHits: 0}
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fsHits++
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) getHits(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hits: ", cfg.fsHits)
}

func (cfg *apiConfig) resetHits(w http.ResponseWriter, r *http.Request) {
	cfg.fsHits = 0
}

func main() {
	mux := http.NewServeMux()
	apiCfg := newApiConfig()
	fsHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/*", apiCfg.middlewareMetricsInc(fsHandler))
	mux.HandleFunc("GET /reset", apiCfg.resetHits)
	mux.HandleFunc("GET /metrics", apiCfg.getHits)
	mux.HandleFunc("GET /healthz", checkHealth)
	corsMux := middlewareCors(mux)
	server := http.Server{Addr: "localhost:8080", Handler: corsMux}
	err := server.ListenAndServe()
	fmt.Println(err)
}
