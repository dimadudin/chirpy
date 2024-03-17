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
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	respParams := responseParameters{Id: newUser.Id, Email: newUser.Email, IsChirpyRed: newUser.IsChirpyRed}
	RespondWithJSON(w, http.StatusCreated, respParams)
}

func (cfg *Config) ApiLogin(w http.ResponseWriter, r *http.Request) {
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

	refreshExpTime := time.Hour * 24 * 60
	refreshClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		Subject:   strconv.Itoa(user.Id),
		ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(refreshExpTime)},
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
	}

	accessExpTime := time.Hour
	accessClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		Subject:   strconv.Itoa(user.Id),
		ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(accessExpTime)},
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	_, err = cfg.db.CreateToken(refreshTokenStr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type responseParameters struct {
		Id           int    `json:"id"`
		Email        string `json:"email"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	respParams := responseParameters{
		Id:           user.Id,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		AccessToken:  accessTokenStr,
		RefreshToken: refreshTokenStr,
	}
	RespondWithJSON(w, http.StatusOK, respParams)
}

func (cfg *Config) ApiUpdateUser(w http.ResponseWriter, r *http.Request) {
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
		Id          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}
	respParams := responseParameters{Id: user.Id, Email: user.Email, IsChirpyRed: user.IsChirpyRed}
	RespondWithJSON(w, http.StatusOK, respParams)
}

func (cfg *Config) ApiUpgradeUser(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		RespondWithError(w, http.StatusUnauthorized, errors.New("no auth").Error())
		return
	}

	apiKey := strings.TrimPrefix(auth, "ApiKey ")
	if apiKey != cfg.polkaApiKey {
		RespondWithError(w, http.StatusUnauthorized, errors.New("wrong api key").Error())
		return
	}

	type requestParameters struct {
		Data struct {
			UserId int `json:"user_id"`
		} `json:"data"`
		Event string `json:"event"`
	}
	rqParams := requestParameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rqParams)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if rqParams.Event != "user.upgraded" {
		RespondWithJSON(w, http.StatusOK, struct{}{})
		return
	}

	_, err = cfg.db.GetUserByID(rqParams.Data.UserId)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	_, err = cfg.db.UpgradeUser(rqParams.Data.UserId)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, struct{}{})
}
