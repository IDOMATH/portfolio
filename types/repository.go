package types

import (
	"fmt"
	"github.com/IDOMATH/portfolio/handlers"
	"github.com/IDOMATH/portfolio/middleware"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/util"
	"github.com/IDOMATH/session"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var regexNumber = regexp.MustCompile(`\d`)

type Repository struct {
	Session  session.Store
	BH       *handlers.BlogHandler
	AH       *handlers.AuthHandler
	GH       *handlers.GuestbookHandler
	FH       *handlers.FitnessHandler
	urlIndex int
}

func NewRepo() *Repository {
	return &Repository{}
}

func (repo *Repository) Route(w http.ResponseWriter, r *http.Request) {
	repo.urlIndex = 1
	url := strings.Split(r.URL.Path, "/")
	switch url[repo.urlIndex] {
	case "":
		handlers.HandleHome(w, r)
	case "contact":
		handlers.HandleContact(w, r)
	case "blog":
		repo.routeBlog(w, r)
	case "pic":
		HandlePic(w, r)
	case "resume":
		handlers.HandleGetResume(w, r)
	case "guestbook":
		repo.routeGuestbook(w, r)
	case "user":
		repo.AH.HandleUserSignUp(w, r)
	case "fitness":
		repo.FH.HandleGetFitness(w, r)
	case "clicked":
		handleClicked(w, r)
	case "admin":
		middleware.Authentication(repo.routeAdmin, repo, w, r)(w, r)

	default:
		handle404(w, r)
	}
}

// routeBlog handles the url segment /blog
func (repo *Repository) routeBlog(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > repo.urlIndex {
		repo.urlIndex++
		segment := url[repo.urlIndex]
		switch {
		// This is /blog/{id}
		case regexNumber.MatchString(segment):
			repo.BH.HandleGetBlogById(w, r)
		}
	}
	repo.BH.HandleBlog(w, r)
}

func (repo *Repository) routeGuestbook(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > repo.urlIndex {
		repo.urlIndex++
		switch segment := url[repo.urlIndex]; segment {
		case "sign":
			repo.GH.HandlePostGuestbookSignature(w, r)
		}
	}
	repo.GH.HandleGetApprovedGuestbookSignatures(w, r)
}

func (repo *Repository) routeAdmin(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > repo.urlIndex {
		repo.urlIndex++
		switch segment := url[repo.urlIndex]; segment {
		case "guestbook":
			middleware.Authentication(repo.routeAdminGuestbook, repo, w, r)
		case "blog":
			repo.routeAdminBlog(w, r)
		case "fitness":
			repo.FH.HandlePostFitness(w, r)
		}
	}
}

func (repo *Repository) routeAdminBlog(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > repo.urlIndex {
		repo.urlIndex++
		switch segment := url[repo.urlIndex]; segment {
		case "new":
			repo.BH.HandleNewBlog(w, r)
		default:
			//TODO: Implement a dashboard to get all blogs to edit
		}
	}
}

func (repo *Repository) routeAdminGuestbook(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > repo.urlIndex {
		repo.urlIndex++
		switch segment := url[repo.urlIndex]; segment {
		case "approve":
			repo.GH.HandleApproveGuestbookSignature(w, r)
		case "deny":
			repo.GH.HandleDenyGuestbookSignature(w, r)
		}
	}
	repo.GH.HandleGetAllGuestbookSignature(w, r)
}

func handle404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	render.Template(w, r, "error-404.go.html", &TemplateData{PageTitle: "Not Found"})
}

func handleClicked(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "clicked.go.html", &TemplateData{})
}

func HandlePic(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := render.Template(w, r, "upload-pic.go.html",
			&TemplateData{PageTitle: "Pic"})
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
