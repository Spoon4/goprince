package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "This is the RESTful api")
}

func luckyHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.URL.Query()["name"][0]
	fmt.Fprintf(w, "Hello, %q\nYour lucky number  is: %d", name, generateLuckyNumber(name))
}

func generateLuckyNumber(name string) int {
	return len(name)
}

func main() {

	// print env
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/lucky", luckyHandler)

	http.ListenAndServe(":8080", router)
}
