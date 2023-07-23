package main

import (
	"context"
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofor-little/env"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/smtp"
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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}
	blogHandler := handlers.NewBlogHandler(db.NewBlogStore(client, dbName))

	app := fiber.New(config)

	app.Get("/", HandleHome)

	app.Get("/contact", HandleGetContactForm)
	app.Post("/contact", HandlePostContactForm)

	app.Get("/blog", blogHandler.HandleGetBlogs)
	app.Post("/blog", blogHandler.HandlePostBlog)
	app.Get("/blog/:id", blogHandler.HandleGetBlogById)

	app.Listen(portNumber)
}

func HandleHome(c *fiber.Ctx) error {
	return c.Render("home", fiber.Map{"PageTitle": "Home"})
}

func HandleGetContactForm(c *fiber.Ctx) error {
	return c.Render("new-blog", fiber.Map{"PageTitle": "Contact"})
}

type ContactDetails struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func HandlePostContactForm(c *fiber.Ctx) error {
	details := ContactDetails{}

	err := c.BodyParser(&details)
	if err != nil {
		return err
	}

	from, err := env.MustGet("EMAIL")
	if err != nil {
		return err
	}

	password, err := env.MustGet("PASSWORD")
	if err != nil {
		return err
	}

	// TODO: some sort of validation for the email
	to := []string{
		details.Email,
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte(details.Message)

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return err
	}

	fmt.Sprintf("details: %x", details)
	//TODO: something about rendering a template using the details coming in
	//TODO: that also means making a new template for this to render
	return c.Render("contact-submitted", fiber.Map{"PageTitle": "Contact Submitter",
		"ContactDetails": details})
}
