package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"siteapi/app/models"
	"siteapi/app/utils"
	"strconv"

	"github.com/gorilla/mux"
)

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to Site API")
	}
}

func GetSites() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sites, err := models.GetSites()
		if err != nil {
			log.Printf("Cannot get sites, err=%v \n", err)
			utils.SendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sites)
	}
}

func GetSiteById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		if _, err := strconv.Atoi(id); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		sites, err := models.GetSiteById(id)
		if err != nil {
			log.Printf("Cannot get site, err=%v \n", err)
			utils.SendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(sites)
	}
}

func CreateSite() http.HandlerFunc {
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

		err = models.CreateSite(site)
		if err != nil {
			log.Printf("Cannot save site in DB. err=%v \n", err)
			utils.SendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		utils.SendResponse(w, r, "Successfully added", http.StatusOK)
	}
}

func UpdateSite() http.HandlerFunc {
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

		err = models.UpdateSite(site)
		if err != nil {
			log.Printf("Cannot update site in DB. err=%v \n", err)
			utils.SendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		utils.SendResponse(w, r, "Successfully updated", http.StatusOK)
	}
}

func DeleteSite() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		if _, err := strconv.Atoi(id); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		err := models.DeleteSite(id)
		if err != nil {
			log.Printf("Cannot removed site, err=%v \n", err)
			utils.SendResponse(w, r, nil, http.StatusInternalServerError)
			return
		}
		utils.SendResponse(w, r, "Successfully removed", http.StatusOK)
	}
}

func Monitoring() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		models.MonitoringSite()
	}
}
