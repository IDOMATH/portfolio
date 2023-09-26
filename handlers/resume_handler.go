package handlers

import (
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"net/http"
)

func HandleGetResume(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "resume.go.html",
		&types.TemplateData{
			PageTitle: "Resume",
		})
}
