package main

import (
	. "authentication-service/config"
	"authentication-service/controllers"
	"authentication-service/dao"
	"authentication-service/middlewares"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
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
	router.Handle("/resource", negroni.New(
		negroni.HandlerFunc(middlewares.ValidateMiddleware),
		negroni.Wrap(http.HandlerFunc(controllers.ProtectedHandler)),
	))

	if err := http.ListenAndServe("localhost:3000", router); err != nil {
		log.Fatal(err)
	}
}