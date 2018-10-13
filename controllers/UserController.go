package controllers

import (
	"authentication-service/dao"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)
var myDao = dao.UserDAO{}


func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := myDao.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := myDao.GetUser(params["id"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid User Id")
		return
	}
	respondWithJson(w, http.StatusOK, user)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}