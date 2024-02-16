package middleware

import (
	"github.com/IDOMATH/portfolio/types"
	"log"
	"net/http"
)

func Authentication(next http.HandlerFunc, repo *types.Repository, w http.ResponseWriter, r *http.Request) http.HandlerFunc {
	// TODO: Make this actually check authentication of the user
	log.Println("Not Authenticated")
	return func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("authToken")
		isTokenMatch, err := repo.SS.CheckSessionToken(headerToken)
		if err != nil {
			// TODO: Render some error page (not 401)
			return
		}
		if isTokenMatch {
			next(w, r)
		}
		// TODO: Render 401
	}
}
