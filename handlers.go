package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (cfg *Config) ApiPostChirp(w http.ResponseWriter, r *http.Request) {
	type requestParameters struct {
		Body string `json:"body"`
	}
	rqParams := requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rqParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(rqParams.Body) > 140 || len(rqParams.Body) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp")
		return
	}
	censored := CensorChirp(rqParams.Body)
	newChirp, err := cfg.db.CreateChirp(censored)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusCreated, newChirp)
}

func (cfg *Config) ApiGetChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *Config) ApiGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirps, err := cfg.db.GetChirpByID(chirpID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, chirps)
}
