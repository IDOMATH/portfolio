package main

import (
	"github.com/IDOMATH/portfolio/config"
	"github.com/IDOMATH/portfolio/render"
	"net/http"
)

const portNumber = ":8080"

var app config.AppConfig

// main is the entry point to the application
func main() {
	render.NewRenderer(&app)
	getHandlers()
	http.ListenAndServe(":8080", nil)
}

func getHandlers() {
	http.HandleFunc("/", Home)
}

func Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl")
}
