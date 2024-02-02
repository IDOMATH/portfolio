package middleware

import (
	"log"
	"net/http"
)

func Authentication(next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	// TODO: Make this actually check authentication of the user
	log.Println("Not Authenticated")
	next(w, r)
}
