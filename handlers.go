package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
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
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	censor := "****"
	words := strings.Split(rqParams.Body, " ")
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = censor
		}
	}
	type responseParameters struct {
		CleanedBody string `json:"cleaned_body"`
	}
	respParams := responseParameters{CleanedBody: strings.Join(words, " ")}
	RespondWithJSON(w, http.StatusOK, respParams)
}
