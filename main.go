package main

import (
	"fmt"
	"github.com/IDOMATH/portfolio/handlers"
	"github.com/IDOMATH/portfolio/render"
	"github.com/IDOMATH/portfolio/types"
	"github.com/IDOMATH/portfolio/util"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofor-little/env"
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
	//client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//blogHandler := handlers.NewBlogHandler(db.NewBlogStore(client, dbName))

	//app := fiber.New(config)

	//app.Get("/", HandleHome)
	http.HandleFunc("/", HandleHome)

	http.HandleFunc("/contact", handlers.HandleContact)
	//app.Get("/contact", HandleGetContactForm)
	//app.Post("/contact", HandlePostContactForm)
	//
	//app.Get("/blog", blogHandler.HandleGetBlogs)
	//app.Post("/blog", blogHandler.HandlePostBlog)
	//app.Get("/blog/:id", blogHandler.HandleGetBlogById)
	//
	//app.Get("/pic", handleGetPic)
	//app.Post("/pic", HandleFileUpload)
	//
	http.HandleFunc("/resume", handlers.HandleGetResume)

	//app.Listen(portNumber)
	http.ListenAndServe(portNumber, nil)
}

func handleGetPic(c *fiber.Ctx) error {
	return c.Render("upload-pic", fiber.Map{"PageTitle": "Upload"}, "layouts/base")
}

func WriteFile(fileName string, file []byte) error {
	return os.WriteFile(fmt.Sprintf("./uploads/%s", fileName), file, 0644)
}

func HandleHome(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "home.go.html",
		&types.TemplateData{PageTitle: "Home"})
	if err != nil {
		util.WriteError(w, http.StatusInternalServerError, err)
	}
}

// TODO: Finish handling file uploads
func HandleFileUpload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	fmt.Println(form.File["file"])
	return nil
}
