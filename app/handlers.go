package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"siteapi/app/models"
	"strconv"

	"github.com/gorilla/mux"
)

func (a *App) IndexHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Site API")
	}
}

func (a *App) GetSitesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sites, err := a.DB.GetSites()
		if err != nil {
			log.Printf("Cannot get sites, err=%v \n", err)
			sendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sites)
	}
}

func (a *App) GetSiteByIdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		if _, err := strconv.Atoi(id); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		sites, err := a.DB.GetSiteById(id)
		if err != nil {
			log.Printf("Cannot get site, err=%v \n", err)
			sendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sites)
	}
}

func (a *App) CreateSiteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(err)
		}
		site := models.Site{}
		err = json.Unmarshal(body, &site)
		if err != nil {
			panic(err)
		}

		err = a.DB.CreateSite(site)
		if err != nil {
			log.Printf("Cannot save site in DB. err=%v \n", err)
			sendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		sendResponse(w, r, "Successfully added", http.StatusOK)
	}
}

func (a *App) UpdateSiteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)

		if err != nil {
			panic(err)
		}
		site := models.Site{}
		err = json.Unmarshal(body, &site)
		if err != nil {
			panic(err)
		}

		err = a.DB.UpdateSite(site)
		if err != nil {
			log.Printf("Cannot update site in DB. err=%v \n", err)
			sendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		sendResponse(w, r, "Successfully updated", http.StatusOK)
	}
}

func (a *App) DeleteSiteHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		if _, err := strconv.Atoi(id); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err := a.DB.DeleteSite(id)
		if err != nil {
			log.Printf("Cannot removed site, err=%v \n", err)
			sendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		sendResponse(w, r, "Successfully removed", http.StatusOK)
	}
}

func (a *App) MonitoringHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.DB.MonitoringSite()
	}
}
