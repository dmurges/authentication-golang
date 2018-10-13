package main

import (
	. "authentication-service/config"
	"authentication-service/controllers"
	"authentication-service/dao"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var myDao = dao.UserDAO{}
var config = Config{}


func init() {
	config.Read()

	myDao.Server = config.Server
	myDao.Database = config.Database
	myDao.Connect()
}


func main() {
	InitKeys()
	router := mux.NewRouter()
	router.HandleFunc("/signup", controllers.Signup).Methods("POST")
	router.HandleFunc("/login", controllers.LoginController).Methods("POST")
	router.HandleFunc("/user", controllers.GetAllUsers).Methods("GET")
	router.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/authenticate", CreateTokenEndpoint).Methods("POST")
	router.HandleFunc("/protected", ProtectedEndpoint).Methods("GET")
	router.HandleFunc("/test", ValidateMiddleware(TestEndpoint)).Methods("GET")

	if err := http.ListenAndServe("localhost:3000", router); err != nil {
		log.Fatal(err)
	}
}