package middleware

import (
	"log"
	"net/http"

	"github.com/IDOMATH/portfolio/repository"
)

func Authentication(next http.HandlerFunc, repo *repository.Repository) http.HandlerFunc {
	// TODO: Make this actually check authentication of the user
	log.Println("Not Authenticated")
	return func(w http.ResponseWriter, r *http.Request) {
		headerToken := r.Header.Get("authToken")
		_, found, err := repo.Session.Get(headerToken)
		if err != nil {
			// TODO: Render some error page (not 401)
			w.Write([]byte("Error getting token from session: " + err.Error()))
			return
		}
		if found {
			next(w, r)
		}
		// TODO: Render 401
		w.Write([]byte("Token not found"))
	}
}
