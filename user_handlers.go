package main

import (
	"encoding/json"
	"net/http"

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
	users, err := cfg.db.GetUsers()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, user := range users {
		if user.Email == rqParams.Email {
			err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(rqParams.Password))
			if err != nil {
				RespondWithError(w, http.StatusUnauthorized, err.Error())
				return
			}
			type responseParameters struct {
				Id    int    `json:"id"`
				Email string `json:"email"`
			}
			respParams := responseParameters{Id: user.Id, Email: user.Email}
			RespondWithJSON(w, http.StatusOK, respParams)
			return
		}
	}
	RespondWithError(w, http.StatusInternalServerError, err.Error())
}

func (cfg *Config) ApiUpdateUser(w http.ResponseWriter, r *http.Request) {}
