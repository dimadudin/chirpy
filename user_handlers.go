package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *Config) ApiCreateUser(w http.ResponseWriter, r *http.Request) {
	type requestParameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	rqParams := requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rqParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rqParams.Password), 0)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newUser, err := cfg.db.CreateUser(rqParams.Email, string(hashedPassword))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type responseParameters struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	respParams := responseParameters{Id: newUser.Id, Email: newUser.Email}
	RespondWithJSON(w, http.StatusCreated, respParams)
}

func (cfg *Config) ApiLogin(w http.ResponseWriter, r *http.Request) {
	type requestParameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds string `json:"expires_in_seconds"`
	}
	rqParams := requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rqParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.GetUserByEmail(rqParams.Email)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rqParams.Password))
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	expTime := time.Hour * 24
	if rqParams.ExpiresInSeconds != "" {
		expTime, err := time.ParseDuration(rqParams.ExpiresInSeconds)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		if expTime > time.Hour*24 {
			expTime = time.Hour * 24
		}
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   strconv.Itoa(user.Id),
		ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(expTime)},
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type responseParameters struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
	}
	respParams := responseParameters{Id: user.Id, Email: user.Email, Token: tokenStr}
	RespondWithJSON(w, http.StatusOK, respParams)
}

func (cfg *Config) ApiUpdateUser(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		RespondWithError(w, http.StatusInternalServerError, errors.New("no auth").Error())
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
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	rqParams := requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&rqParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rqParams.Password), 0)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.db.UpdateUser(userID, rqParams.Email, string(hashedPassword))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type responseParameters struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}
	respParams := responseParameters{Id: user.Id, Email: user.Email}
	RespondWithJSON(w, http.StatusOK, respParams)
}
