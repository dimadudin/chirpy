package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func (cfg *Config) ApiRefreshToken(w http.ResponseWriter, r *http.Request) {
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
	if issuer != "chirpy-refresh" {
		RespondWithError(w, http.StatusUnauthorized, errors.New("invalid token issuer").Error())
		return
	}

	db_token, err := cfg.db.GetToken(tokenStr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !db_token.RevokedAt.IsZero() {
		RespondWithError(w, http.StatusUnauthorized, errors.New("this token has been revoked").Error())
		return
	}

	userIDStr, err := token.Claims.GetSubject()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	accessExpTime := time.Hour
	accessClaims := &jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		Subject:   userIDStr,
		ExpiresAt: &jwt.NumericDate{Time: time.Now().UTC().Add(accessExpTime)},
		IssuedAt:  &jwt.NumericDate{Time: time.Now().UTC()},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type responseParameters struct {
		Token string `json:"token"`
	}
	respParams := responseParameters{Token: accessTokenStr}
	RespondWithJSON(w, http.StatusOK, respParams)
}

func (cfg *Config) ApiRevokeToken(w http.ResponseWriter, r *http.Request) {
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
	if issuer != "chirpy-refresh" {
		RespondWithError(w, http.StatusUnauthorized, errors.New("invalid token issuer").Error())
		return
	}

	_, err = cfg.db.RevokeToken(tokenStr)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type responseParameters struct{}
	respParams := responseParameters{}
	RespondWithJSON(w, http.StatusOK, respParams)
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
