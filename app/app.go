package app

import (
	"siteapi/app/database"

	"github.com/gorilla/mux"
)

type App struct {
	Router *mux.Router
	DB     database.SiteDB
}

func New() *App {
	a := &App{
		Router: mux.NewRouter(),
	}

	a.initRoutes()
	return a
}

func (a *App) initRoutes() {
	a.Router.HandleFunc("/", a.IndexHandler()).Methods("GET")
	a.Router.HandleFunc("/api/sites", a.GetSitesHandler()).Methods("GET")
	a.Router.HandleFunc("/api/sites/{id}", a.GetSiteByIdHandler()).Methods("GET")
	a.Router.HandleFunc("/api/sites", a.CreateSiteHandler()).Methods("POST")
	a.Router.HandleFunc("/api/sites", a.UpdateSiteHandler()).Methods("PUT")
	a.Router.HandleFunc("/api/sites/{id}", a.DeleteSiteHandler()).Methods("DELETE")
	a.Router.HandleFunc("/api/monitorings/executions", a.MonitoringHandler())
}
