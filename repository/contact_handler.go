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

func HandleContact(w http.ResponseWriter, r *http.Request) {
	var err error
	if r.Method == "GET" {
		err = render.Template(w, r, "new-blog.go.html", &render.TemplateData{PageTitle: "Contact"})
	}
	if r.Method == "POST" {
		err = PostContactForm(w, r)
	}
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
}

func PostContactForm(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
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
		return err
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

	return err
}
