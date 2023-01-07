package routes

import (
	"siteapi/app/controllers"

	"github.com/gorilla/mux"
)

var RegisterSiteRoutes = func(router *mux.Router) {
	router.HandleFunc("/", controllers.Index()).Methods("GET")
	router.HandleFunc("/api/sites", controllers.GetSites()).Methods("GET")
	router.HandleFunc("/api/sites/{id}", controllers.GetSiteById()).Methods("GET")
	router.HandleFunc("/api/sites", controllers.CreateSite()).Methods("POST")
	router.HandleFunc("/api/sites", controllers.UpdateSite()).Methods("PUT")
	router.HandleFunc("/api/sites/{id}", controllers.DeleteSite()).Methods("DELETE")
	router.HandleFunc("/api/monitorings/executions", controllers.Monitoring())
}
