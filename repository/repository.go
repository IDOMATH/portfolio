package repository

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"

	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/util"
	"github.com/IDOMATH/session"
)

var regexNumber = regexp.MustCompile(`\d`)

type Repository struct {
	Session  session.Store
	BH       *BlogHandler
	AH       *AuthHandler
	GH       *GuestbookHandler
	FH       *FitnessHandler
	urlIndex int
}

func NewRepo() *Repository {
	return &Repository{}
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render.Template(w, r, "error-404.go.html", &render.TemplateData{PageTitle: "Not Found"})
}

func handleClicked(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "clicked.go.html", &render.TemplateData{})
}

func HandlePic(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := render.Template(w, r, "upload-pic.go.html",
			&render.TemplateData{PageTitle: "Pic"})
		if err != nil {
			util.WriteError(w, http.StatusInternalServerError, err)
		}
	case "POST":
		uploadFile(w, r)
	}
}

// TODO: Make this insert the file location into the DB
func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("mediaFile")
	if err != nil {
		fmt.Println("Error retrieving file")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded file: %+v\n", handler.Filename)
	fmt.Printf("File size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", handler.Filename))
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Copy uploaded file to the created file on filesystem
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Successfully uploaded file\n")
}
