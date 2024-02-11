package middleware

import (
	"log"
	"net/http"
)

func Authentication(next http.HandlerFunc, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	// TODO: Make this actually check authentication of the user
	log.Println("Not Authenticated")
	return func(w http.ResponseWriter, r *http.Request) {
		// Something like get auth token from r.Header.Get()
		r.Header.Get("authToken")
		// Then look up that token in the DB and either continue routing or render the 401 page
		next(w, r)
	}
}
