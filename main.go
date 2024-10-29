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
	repo.FH = repository.NewFitnessHandler(*db.NewPostgresFitnessStore(postgresDb.SQL))
	repo.GH = repository.NewGuestbookHandler(*db.NewPostgresGuestbookStore(postgresDb.SQL))

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
	router.HandleFunc("GET /blog/{id}/", repo.BH.HandleGetBlogById)
	router.HandleFunc("GET /new-blog/", repo.BH.HandleNewBlog)
	router.HandleFunc("POST /new-blog/", repo.BH.HandlePostNewBlog)

	router.HandleFunc("GET /fitness/", repo.FH.HandleGetFitness)
	router.HandleFunc("GET /fitness-form/", repo.FH.HandleGetFitnessForm)
	router.HandleFunc("POST /fitness-form/", repo.FH.HandlePostFitnessForm)

	router.HandleFunc("GET /guestbook/", repo.GH.HandleGetGuestbook)
	router.HandleFunc("POST /guestbook", repo.GH.HandlePostGuestbook)
	router.HandleFunc("GET /admin/guestbook", repo.GH.HandleGetGuestbookAdmin)
	router.HandleFunc("POST /admin/guestbook/approve", repo.GH.HandleApproveGuestbookSignature)
	router.HandleFunc("POST /admin/guestbook/deny", repo.GH.HandleDenyGuestbookSignature)

	repo.AH = repository.NewAuthHandler(db.NewUserStore(client, mongoDbName))

	repo.Session = memStore

	fmt.Println("Starting server on port ", portNumber)
	log.Fatal(server.ListenAndServe())
}
