package config

import (
	. "authentication-service/models"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"log"
	"net/http"
)

type JwtToken struct {
	Token string `json:"token"`
}

type Exception struct {
	Message string `json:"message"`
}

var VerifyKey, SignKey []byte

const (
	privKeyPath = "app.rsa"
	pubKeyPath = "app.rsa.pub"
)

func InitKeys() {
	var err error

	SignKey, err = ioutil.ReadFile(privKeyPath)
	if err != nil {
		log.Fatal("Error reading private key")
		return
	}
	VerifyKey, err = ioutil.ReadFile(pubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key")
		return
	}
}



func CreateTokenEndpoint(w http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
	})
	tokenString, error := token.SignedString([]byte("secret"))
	if error != nil {
		fmt.Println(error)
	}
	json.NewEncoder(w).Encode(JwtToken{Token:tokenString})

}

func ProtectedEndpoint(w http.ResponseWriter, req *http.Request) {
	params := req.URL.Query()
	token, _ := jwt.Parse(params["token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return []byte("secret"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var user User
		mapstructure.Decode(claims, &user)
		json.NewEncoder(w).Encode(user)
	} else {
		json.NewEncoder(w).Encode(Exception{Message:"Invalid Authorization token"})
	}
}

func TestEndpoint(w http.ResponseWriter, req *http.Request) {
	decoded := context.Get(req, "decoded")
	var user User
	mapstructure.Decode(decoded.(jwt.MapClaims), &user)
	json.NewEncoder(w).Encode(user)
}
