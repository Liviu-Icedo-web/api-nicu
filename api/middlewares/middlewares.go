package middlewares

import (
	"errors"
	"fmt"
	"net/http"

	"api-nicu/api/auth"
	"api-nicu/api/responses"
)

func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("*** SetMiddlewareJSON", r, "\n")
		if r.Method == "OPTIONS" {
			return
		}
		next(w, r)
	}
}

func SetMiddlewareAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("**** SetMiddlewareAuthentication", r, "\n")
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}
