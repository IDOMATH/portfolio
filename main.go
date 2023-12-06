package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/handlers"
	"github.com/IDOMATH/portfolio/middleware"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofor-little/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
)

const portNumber = ":8080"
const dbUri = "mongodb://localhost:27017"
const dbName = "portfolio"
const blogCollection = "blog"
const templatesLocation = "./templates"

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
	Views: html.New(templatesLocation, ".go.html"),
}

// main is the entry point to the application
func main() {
	// Get environment variables
	err := env.Load("dev.env")
	if err != nil {
		log.Fatal(err)
	}

	// Set config details based whether the values are in the .env
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	var postgresDb *sql.DB
	//connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	//postgresDb, err := db.ConnectSQL(connectionString)

	blogHandler := handlers.NewBlogHandler(db.NewBlogStore(client, dbName))
	userHandler := handlers.NewUserHandler(db.NewUserStore(client, dbName))
	guestbookHandler := handlers.NewGuestbookHandler(*db.NewPostgresGuestbookStore(postgresDb))

	//app.Get("/", HandleHome)
	http.HandleFunc("/", middleware.Authentication(handlers.HandleHome))

	http.HandleFunc("/contact", handlers.HandleContact)
	//
	http.HandleFunc("/blog", blogHandler.HandleGetBlogs(context.Background()))
	//app.Post("/blog", blogHandler.HandlePostBlog)
	//app.Get("/blog/:id", blogHandler.HandleGetBlogById)
	//
	//app.Get("/pic", handleGetPic)
	http.HandleFunc("/pic", HandlePic)
	//
	http.HandleFunc("/resume", handlers.HandleGetResume)

	http.HandleFunc("/user", userHandler.HandlePostUser)

	http.HandleFunc("/guestbook", guestbookHandler.HandleGetGuestbookSignatures())

	//app.Listen(portNumber)
	http.ListenAndServe(portNumber, nil)
}

func handleGetPic(c *fiber.Ctx) error {
	//types.BlogCard{secretF}
	return c.Render("upload-pic", fiber.Map{"PageTitle": "Upload"}, "layouts/base")
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
