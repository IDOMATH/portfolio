package repository

import (
	"net/http"

	"github.com/IDOMATH/portfolio/render"
)

func HandleGetResume(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "resume.go.html",
		&render.TemplateData{
			PageTitle: "Resume",
		})
}
