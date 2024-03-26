package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/leandrofars/oktopus/internal/api/auth"
	"github.com/leandrofars/oktopus/internal/db"
)

func (a *Api) retrieveUsers(w http.ResponseWriter, r *http.Request) {
	users, err := a.db.FindAllUsers()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, x := range users {
		delete(x, "password")
	}

	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		log.Println(err)
	}
	return
}

func (a *Api) registerUser(w http.ResponseWriter, r *http.Request) {

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	email, err := auth.ValidateToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Check if user which is requesting creation has the necessary privileges
	rUser, err := a.db.FindUser(email)
	if rUser.Level != AdminUser {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var user db.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.Level = NormalUser

	if err := user.HashPassword(user.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.db.RegisterUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *Api) registerAdminUser(w http.ResponseWriter, r *http.Request) {

	var user db.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := a.db.FindAllUsers()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	adminExists := adminUserExists(users)
	if adminExists {
		log.Println("There might exist only one admin")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("There might exist only one admin")
		return
	}

	user.Level = AdminUser

	if err := user.HashPassword(user.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := a.db.RegisterUser(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func adminUserExists(users []map[string]interface{}) bool {
	for _, x := range users {
		if x["level"].(int32) == AdminUser {
			log.Println("Admin exists")
			return true
		}
	}
	return false
}

func (a *Api) adminUserExists(w http.ResponseWriter, r *http.Request) {

	users, err := a.db.FindAllUsers()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	adminExits := adminUserExists(users)
	json.NewEncoder(w).Encode(adminExits)
	return
}

type TokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Api) generateToken(w http.ResponseWriter, r *http.Request) {
	var tokenReq TokenRequest

	err := json.NewDecoder(r.Body).Decode(&tokenReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := a.db.FindUser(tokenReq.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Invalid Credentials")
		return
	}

	credentialError := user.CheckPassword(tokenReq.Password)
	if credentialError != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Invalid Credentials")
		return
	}

	token, err := auth.GenerateJWT(user.Email, user.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
	return
}
