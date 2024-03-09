package main

import (
	"fmt"
	"io"
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

func ApiValidateChirp(w http.ResponseWriter, r *http.Request) {
	dat, err := io.ReadAll(r.Body)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(dat) > 140 || len(dat) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp")
		return
	}
	resp := struct {
		Valid bool `json:"valid"`
	}{Valid: true}
	RespondWithJSON(w, http.StatusOK, resp)
}
