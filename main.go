package main

import (
	"fmt"
	"net/http"
)

// Handler functions
func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home")
}

func main() {
	// Register handlers
	http.HandleFunc("/", home)

	// Start a web server
	http.ListenAndServe(":8080", nil)
}
