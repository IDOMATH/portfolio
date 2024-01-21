package middleware

import (
	"log"
	"net/http"
)

func Authentication(next http.HandlerFunc) http.HandlerFunc {
	// TODO: Make this actually check authentication of the user
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("URL: ", r.URL)
		next(w, r)
	}
}
