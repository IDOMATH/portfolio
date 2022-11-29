package main

import (
	"fmt"
	"github.com/IDOMATH/portfolio/router"
	"log"
	"net/http"
	"strconv"
)

const PORT = "8080"

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome home.")
}

func say(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "You asked me to say: %s", router.UrlParam(req, "Message"))
}

func factors(w http.ResponseWriter, req *http.Request) {
	numberString := router.UrlParam(req, "Number")
	number, err := strconv.Atoi(numberString)
	if err != nil {
		fmt.Fprintln(w, "You did not give me an integer.")
		return
	}
	fmt.Fprintf(w, "The factors of %d are:\n", number)
	lower := 1
	higher := number
	for lower <= higher {
		if number%lower == 0 {
			fmt.Fprintf(w, "%s x %s\n", strconv.Itoa(lower), strconv.Itoa(number/lower))
		}
		lower++
		higher = number / lower
	}
}

func main() {
	// Our router to handle routes
	router := &router.Router{}

	// Add routes
	router.Route("GET", "/", home)
	router.Route(http.MethodGet, `/say/(?P<Message>\w+)`, say)
	router.Route(http.MethodGet, `/factorsOf/(?P<Number>\w+)`, factors)

	log.Printf("Running on port %s", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), router))
}
