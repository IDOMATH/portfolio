package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/repository"
	"github.com/IDOMATH/portfolio/util"
	"github.com/IDOMATH/session/memorystore"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const portNumber = ":8080"
const dbUri = "mongodb://localhost:27017"
const mongoDbName = "portfolio"

// Start the URL at 1, because the leading slash makes entry 0 the empty string ""
var urlIndex = 1

// main is the entry point to the application
func main() {
	fmt.Println("Connecting to mongo")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))
	if err != nil {
		log.Fatal(err)
	}

	dbHost := util.GetEnvValue("DBHOST", "localhost")
	dbPort := util.GetEnvValue("DBPORT", "5432")
	dbName := util.GetEnvValue("DBNAME", "portfolio")
	dbUser := util.GetEnvValue("DBUSER", "postgres")
	dbPass := util.GetEnvValue("DBPASS", "postgres")
	dbSsl := util.GetEnvValue("DBSSL", "disable")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", dbHost, dbPort, dbName, dbUser, dbPass, dbSsl)
	fmt.Println("Connecting to Postgres")
	postgresDb, err := db.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to Postgres")

	memStore := memorystore.New()

	repo := repository.NewRepo()

	router := http.NewServeMux()
	server := http.Server{
		Addr:    portNumber,
		Handler: router,
	}

	router.HandleFunc("GET /", repository.HandleHome)
	router.HandleFunc("GET /contact/", repository.HandleGetContact)
	router.HandleFunc("POST /conact/", repository.HandlePostContact)
	router.HandleFunc("GET /resume/", repository.HandleGetResume)
	router.HandleFunc("GET /blog/", repo.BH.HandleGetBlog)
	router.HandleFunc("GET /blog/{id}", repo.BH.HandleGetBlogById)
	router.HandleFunc("GET /new-blog/", repo.BH.HandleNewBlog)
	router.HandleFunc("POST /new-blog/", repo.BH.HandlePostNewBlog)

	repo.BH = repository.NewBlogHandler(db.NewBlogStore(client, mongoDbName))
	repo.AH = repository.NewAuthHandler(db.NewUserStore(client, mongoDbName))
	repo.GH = repository.NewGuestbookHandler(*db.NewPostgresGuestbookStore(postgresDb.SQL))
	repo.FH = repository.NewFitnessHandler(*db.NewPostgresFitnessStore(postgresDb.SQL))

	repo.Session = memStore

	// Match all requests and route them with our router
	// http.HandleFunc("/", repo.Route)

	fmt.Println("Starting server on port ", portNumber)
	// http.ListenAndServe(portNumber, nil)
	log.Fatal(server.ListenAndServe())
}
