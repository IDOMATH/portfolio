package main

import (
	"context"
	"fmt"
	"github.com/IDOMATH/portfolio/db"
	"github.com/IDOMATH/portfolio/handlers"
	"github.com/IDOMATH/portfolio/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
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

	repo := types.NewRepo()

	repo.BH = handlers.NewBlogHandler(db.NewBlogStore(client, mongoDbName))
	repo.AH = handlers.NewAuthHandler(db.NewUserStore(client, mongoDbName))
	repo.GH = handlers.NewGuestbookHandler(*db.NewPostgresGuestbookStore(postgresDb.SQL))
	repo.FH = handlers.NewFitnessHandler(*db.NewPostgresFitnessStore(postgresDb.SQL))
	repo.SS = db.NewSessionStore(postgresDb.SQL)

	// Match all requests and route them with our router
	http.HandleFunc("/", repo.Route)

	fmt.Println("Starting server on port ", portNumber)
	http.ListenAndServe(portNumber, nil)
}
