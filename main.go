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

var blogHandler *handlers.BlogHandler
var userHandler *handlers.UserHandler
var guestbookHandler *handlers.GuestbookHandler
var fitnessHandler *handlers.FitnessHandler

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

	blogHandler = handlers.NewBlogHandler(db.NewBlogStore(client, mongoDbName))
	userHandler = handlers.NewUserHandler(db.NewUserStore(client, mongoDbName))
	guestbookHandler = handlers.NewGuestbookHandler(*db.NewPostgresGuestbookStore(postgresDb.SQL))
	fitnessHandler = handlers.NewFitnessHandler(*db.NewPostgresFitnessStore(postgresDb.SQL))

	// Match all requests and route them with our router
	http.HandleFunc("/", Route)

	fmt.Println("Starting server on port ", portNumber)
	http.ListenAndServe(portNumber, nil)
}

func Route(w http.ResponseWriter, r *http.Request) {
	urlIndex = 1
	url := strings.Split(r.URL.Path, "/")
	switch url[urlIndex] {
	case "":
		handlers.HandleHome(w, r)
	case "contact":
		handlers.HandleContact(w, r)
	case "blog":
		routeBlog(w, r)
	case "pic":
		HandlePic(w, r)
	case "resume":
		handlers.HandleGetResume(w, r)
	case "guestbook":
		routeGuestbook(w, r)
	case "user":
		userHandler.HandlePostUser(w, r)
	case "fitness":
		fitnessHandler.HandleGetFitness(w, r)
	case "clicked":
		handleClicked(w, r)
	case "admin":
		middleware.Authentication(routeAdmin, w, r)

	default:
		handle404(w, r)
	}
}

// routeBlog handles the url segment /blog
func routeBlog(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > urlIndex {
		urlIndex++
		segment := url[urlIndex]
		switch {
		// This is /blog/{id}
		case regexNumber.MatchString(segment):
			blogHandler.HandleGetBlogById(w, r)
		}
	}
	blogHandler.HandleBlog(w, r)
}

func routeGuestbook(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > urlIndex {
		urlIndex++
		switch segment := url[urlIndex]; segment {
		case "sign":
			guestbookHandler.HandlePostGuestbookSignature(w, r)
		}
	}
	guestbookHandler.HandleGetApprovedGuestbookSignatures(w, r)
}

func routeAdmin(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > urlIndex {
		urlIndex++
		switch segment := url[urlIndex]; segment {
		case "guestbook":
			middleware.Authentication(routeAdminGuestbook, w, r)
		case "blog":
			blogHandler.HandleNewBlog(w, r)
		case "fitness":
			fitnessHandler.HandlePostFitness(w, r)
		}
	}
}

func routeAdminGuestbook(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.URL.Path, "/")
	if len(url)-1 > urlIndex {
		urlIndex++
		switch segment := url[urlIndex]; segment {
		case "approve":
			guestbookHandler.HandleApproveGuestbookSignature(w, r)
		case "deny":
			guestbookHandler.HandleDenyGuestbookSignature(w, r)
		}
	}
	guestbookHandler.HandleGetAllGuestbookSignature(w, r)
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
