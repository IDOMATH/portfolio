package middleware

import (
	"log"
	"net/http"
)

func Authentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("URL: ", r.URL)
		next(w, r)
	}
}
