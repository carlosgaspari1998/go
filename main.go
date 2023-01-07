package main

import (
	"log"
	"net/http"
	"siteapi/app/routes"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	routes.RegisterSiteRoutes(r)
	http.Handle("/", r)
	log.Println("App running..")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
