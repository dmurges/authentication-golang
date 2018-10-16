package controllers

import (
	"authentication-service/config"
	. "authentication-service/models"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)



type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

type Token struct {
	Token 	string    `json:"token"`
}

type Response struct {
	Data	string	`json:"data"`
}


func LoginController(w http.ResponseWriter, r *http.Request) {
	var credentials Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user, err := myDao.GetUserByEmail(credentials.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid User Id")
		return
	}

	if !comparePasswords(user.Password, []byte(credentials.Password)) {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	claims := JwtClaims{
		user.Username,
		jwt.StandardClaims{
			Id: user.ID.Hex(),
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte(config.SignKey))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	response := Token{token}
	JsonResponse(response, w)
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request){

	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)

}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println("Password is incorrect", err)
		return false
	}
	return true
}


func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err :=  json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}