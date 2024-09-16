package repository

import (
	"net/http"

	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/util"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "home.go.html",
		&render.TemplateData{PageTitle: "Home"})
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
}
