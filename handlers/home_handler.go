package handlers

import (
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"net/http"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "home.go.html",
		&types.TemplateData{PageTitle: "Home"})
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
}
