package main

import (
	"log"
	"net/http"
	"siteapi/app"
	"siteapi/app/database"
)

func main() {
	app := app.New()
	app.DB = &database.DB{}
	err := app.DB.Open()

	if err != nil {
		panic(err)
	}

	defer app.DB.Close()

	http.HandleFunc("/", app.Router.ServeHTTP)

	log.Println("App running..")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
