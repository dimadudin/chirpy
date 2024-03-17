package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *Config) ApiPostChirp(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		RespondWithError(w, http.StatusUnauthorized, errors.New("no auth").Error())
		return
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		})
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if issuer != "chirpy-access" {
		RespondWithError(w, http.StatusUnauthorized, errors.New("invalid token issuer").Error())
		return
	}

	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type requestParameters struct {
		Body string `json:"body"`
	}
	rqParams := requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&rqParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(rqParams.Body) > 140 || len(rqParams.Body) == 0 {
		RespondWithError(w, http.StatusBadRequest, "Invalid chirp")
		return
	}
	censored := CensorChirp(rqParams.Body)
	newChirp, err := cfg.db.CreateChirp(censored, userID)
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
	chirp, err := cfg.db.GetChirpByID(chirpID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	RespondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *Config) ApiDeleteChirpByID(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		RespondWithError(w, http.StatusUnauthorized, errors.New("no auth").Error())
		return
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		})
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if issuer != "chirpy-access" {
		RespondWithError(w, http.StatusUnauthorized, errors.New("invalid token issuer").Error())
		return
	}

	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := cfg.db.GetChirpByID(chirpID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if chirp.AuthorId != userID {
		RespondWithError(w, http.StatusForbidden, errors.New("chirp deletion forbidden").Error())
		return
	}

	cfg.db.DeleteChirp(chirp.Id)

	RespondWithJSON(w, http.StatusOK, chirp)
}
