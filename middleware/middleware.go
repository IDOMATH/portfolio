package middleware

import (
	"log"
	"net/http"
)

func Authentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Not Authenticated")
		next(w, r)
	}
}
