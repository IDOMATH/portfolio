package main

import (
	"context"
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/handlers"
	"github.com/IDOMATH/portfolio/middleware"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const portNumber = ":8080"
const dbUri = "mongodb://localhost:27017"
const mongoDbName = "portfolio"

var regexNumber = regexp.MustCompile(`\d`)

// Start the URL at 1, because the leading slash makes entry 0 the empty string ""
var urlIndex = 1

// main is the entry point to the application
func main() {
	fmt.Println("Connecting to mongo")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	dbHost := "localhost"
	dbPort := "5432"
	dbName := "portfolio"
	dbUser := "postgres"
	dbPass := "postgres"
	dbSSL := "disable"
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSSL)
	fmt.Println("Connecting to Postgres")
	postgresDb, err := db.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Postgres")

	repo := NewRepo()

	repo.BH = handlers.NewBlogHandler(db.NewBlogStore(client, mongoDbName))
	repo.UH = handlers.NewUserHandler(db.NewUserStore(client, mongoDbName))
	repo.GH = handlers.NewGuestbookHandler(*db.NewPostgresGuestbookStore(postgresDb.SQL))
	repo.FH = handlers.NewFitnessHandler(*db.NewPostgresFitnessStore(postgresDb.SQL))

	// Match all requests and route them with our router
	http.HandleFunc("/", repo.Route)

	fmt.Println("Starting server on port ", portNumber)
	http.ListenAndServe(portNumber, nil)
}

type Repository struct {
	Session map[string]string
	BH      *handlers.BlogHandler
	UH      *handlers.UserHandler
	GH      *handlers.GuestbookHandler
	FH      *handlers.FitnessHandler
}

func NewRepo() *Repository {
	return &Repository{}
}

func (repo *Repository) Route(w http.ResponseWriter, r *http.Request) {
	urlIndex = 1
	url := strings.Split(r.URL.Path, "/")
	switch url[urlIndex] {
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
		repo.UH.HandlePostUser(w, r)
	case "fitness":
		repo.FH.HandleGetFitness(w, r)
	case "clicked":
		handleClicked(w, r)
	case "admin":
		middleware.Authentication(repo.routeAdmin, w, r)

	default:
		handle404(w, r)
	}
}

// routeBlog handles the url segment /blog
func (repo *Repository) routeBlog(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > urlIndex {
		urlIndex++
		segment := url[urlIndex]
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
	if len(url)-1 > urlIndex {
		urlIndex++
		switch segment := url[urlIndex]; segment {
		case "sign":
			repo.GH.HandlePostGuestbookSignature(w, r)
		}
	}
	repo.GH.HandleGetApprovedGuestbookSignatures(w, r)
}

func (repo *Repository) routeAdmin(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > urlIndex {
		urlIndex++
		switch segment := url[urlIndex]; segment {
		case "guestbook":
			middleware.Authentication(repo.routeAdminGuestbook, w, r)
		case "blog":
			repo.routeAdminBlog(w, r)
		case "fitness":
			repo.FH.HandlePostFitness(w, r)
		}
	}
}

func (repo *Repository) routeAdminBlog(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > urlIndex {
		urlIndex++
		switch segment := url[urlIndex]; segment {
		case "new":
			repo.BH.HandleNewBlog(w, r)
		default:
			//TODO: Implement a dashboard to get all blogs to edit
		}
	}
}

func (repo *Repository) routeAdminGuestbook(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > urlIndex {
		urlIndex++
		switch segment := url[urlIndex]; segment {
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
	render.Template(w, r, "error-404.go.html", &types.TemplateData{PageTitle: "Not Found"})
}

func handleClicked(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "clicked.go.html", &types.TemplateData{})
}

func HandlePic(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		err := render.Template(w, r, "upload-pic.go.html",
			&types.TemplateData{PageTitle: "Pic"})
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
