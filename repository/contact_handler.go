package repository

import (
	"net/http"

	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/util"
)

type ContactDetails struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func HandleGetContact(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "new-blog.go.html", &render.TemplateData{PageTitle: "Contact"})
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
}

func HandlePostContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	email := r.Form.Get("email")
	subject := r.Form.Get("subject")
	message := r.Form.Get("message")

	bools := make(map[string]bool)

	if !util.IsValidEmail(email) {
		bools["submitted_successfully"] = false
		render.Template(w, r, "contact-submitted.go.html",
			&render.TemplateData{
				PageTitle: "Contact",
				BoolMap:   bools,
			})
	}
	err = util.SendEmail(email, message)
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	bools["submitted_successfully"] = true

	objects := make(map[string]interface{})
	contactDetails := ContactDetails{
		Email:   email,
		Subject: subject,
		Message: message,
	}
	objects["contact_details"] = contactDetails

	err = render.Template(w, r, "contact-submitted.go.html",
		&render.TemplateData{
			PageTitle: "Contact",
			BoolMap:   bools,
			ObjectMap: objects,
		})

	util.WriteError(w, http.StatusInternalServerError, err)
}
