package middleware

import (
	"fmt"
	"log"
	"net/http"
)

func Authentication(next http.HandlerFunc, w http.ResponseWriter, r *http.Request) {
	fmt.Println("authentication")
	// TODO: Make this actually check authentication of the user
	log.Println("Not Authenticated")
	next(w, r)
}
