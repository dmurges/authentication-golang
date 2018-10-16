package middlewares

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"net/http"
	"strings"
)

type Exception struct {
	Message string `json:"message"`
}

//func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
//
//	//validate token
//	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error){
//		return config.VerifyKey, nil
//	})
//
//	if err == nil {
//
//		if token.Valid{
//			next(w, r)
//		} else {
//			w.WriteHeader(http.StatusUnauthorized)
//			fmt.Fprint(w, "Token is not valid")
//		}
//	} else {
//		w.WriteHeader(http.StatusUnauthorized)
//		fmt.Fprint(w, "Unauthorised access to this resource")
//	}
//
//}

func ValidateMiddleware(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	authorizationHeader := req.Header.Get("authorization")
	if authorizationHeader != "" {
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) == 2 {
			token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return []byte("secret"), nil
			})
			if error != nil {
				json.NewEncoder(w).Encode(Exception{Message: error.Error()})
				return
			}
			if token.Valid {
				context.Set(req, "decoded", token.Claims)
				next(w, req)
			} else {
				json.NewEncoder(w).Encode(Exception{Message:"Invalid authorization token"})
			}
		}
	} else {
		json.NewEncoder(w).Encode(Exception{Message:"An authorization header is required"})
	}
}