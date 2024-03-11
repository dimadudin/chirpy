package main

import (
	"fmt"
	"net/http"
)

func (cfg *Config) AdminGetHitCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	body := `
<html>
<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
</body>
</html>
`
	fmt.Fprintf(w, body, cfg.GetHitCount())
}

func (cfg *Config) ApiResetHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.ResetHits()
	fmt.Fprint(w, `Hits have been reset`)
}

func ApiCheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `OK`)
}
