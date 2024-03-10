package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-type", "application/json")
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(dat)
	if err != nil {
		log.Printf("Error writing JSON to header %s:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	resp := struct {
		Error string `json:"error"`
	}{Error: msg}
	RespondWithJSON(w, code, resp)
}

func CensorChirp(chirp string) string {
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	censor := "****"
	words := strings.Split(chirp, " ")
	for i, word := range words {
		if slices.Contains(badWords, strings.ToLower(word)) {
			words[i] = censor
		}
	}
	return strings.Join(words, " ")
}
