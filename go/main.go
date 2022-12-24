package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Site struct {
	ID                  int
	URL                 string
	Frequency           float64
	LastExecutionDate   *time.Time `json:"last_execution_date"`
	Sucess              *bool
	ResponseTime        *float64   `json:"response_time"`
	ResponseAverageTime *float64   `json:"response_average_time"`
	CreationDate        time.Time  `json:"creation_date"`
	LastUpdatedDate     *time.Time `json:"last_updated_date"`
}

func main() {
	log.Println("Server started on: http://localhost:8080")
	router := mux.NewRouter()
	router.HandleFunc("/sites", getSites).Methods("GET")
	router.HandleFunc("/sites/{id}", getSiteById).Methods("GET")
	router.HandleFunc("/sites", createSite).Methods("POST")
	router.HandleFunc("/sites", updateSite).Methods("PUT")
	router.HandleFunc("/sites/{id}", deleteSite).Methods("DELETE")
	router.HandleFunc("/monitorings/executions", monitoring)

	http.ListenAndServe(":8080", router)
}

func getSites(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT id, url, frequency, sucess, response_time, response_average_time, creation_date, last_execution_date, last_updated_date FROM sites")
	if err != nil {
		panic(err.Error())
	}
	site := Site{}
	res := []Site{}
	for selDB.Next() {
		var id int
		var url string
		var frequency float64
		var sucess *bool
		var responseTime, responseAverageTime *float64
		var creationDate time.Time
		var lastExecutionDate, lastUpdatedDate *time.Time
		err = selDB.Scan(&id, &url, &frequency, &sucess, &responseTime, &responseAverageTime, &creationDate, &lastExecutionDate, &lastUpdatedDate)
		if err != nil {
			panic(err.Error())
		}
		site.ID = id
		site.URL = url
		site.Frequency = frequency
		site.Sucess = sucess
		site.ResponseTime = responseTime
		site.ResponseAverageTime = responseAverageTime
		site.CreationDate = creationDate
		site.LastExecutionDate = lastExecutionDate
		site.LastUpdatedDate = lastUpdatedDate

		res = append(res, site)
	}
	json.NewEncoder(w).Encode(res)
	defer db.Close()
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "monitoring"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func getSiteById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	db := dbConn()
	selDB := db.QueryRow("SELECT id, url, frequency, last_execution_date, sucess, response_time, response_average_time, creation_date, last_updated_date FROM sites WHERE id=?", id)
	site := &Site{}

	err := selDB.Scan(&site.ID, &site.URL, &site.Frequency, &site.LastExecutionDate, &site.Sucess, &site.ResponseTime, &site.ResponseAverageTime, &site.CreationDate, &site.LastUpdatedDate)
	if err != nil {
		panic(err.Error())
	}
	json.NewEncoder(w).Encode(site)
	defer db.Close()
}

func createSite(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}
	var site Site
	err = json.Unmarshal(body, &site)
	if err != nil {
		panic(err)
	}

	insDB, err := db.Prepare("INSERT INTO sites(url, frequency, creation_date) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insDB.Exec(site.URL, site.Frequency, site.CreationDate)
	defer db.Close()
}

func updateSite(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}
	var site Site
	err = json.Unmarshal(body, &site)
	if err != nil {
		panic(err)
	}

	time_start := time.Now()
	updDB, err := db.Prepare("UPDATE Sites SET url=?, frequency=?, creation_date=?, last_updated_date=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	updDB.Exec(site.URL, site.Frequency, site.CreationDate, time_start.Format(time.RFC3339), site.ID)
	defer db.Close()
}

func deleteSite(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if _, err := strconv.Atoi(id); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	db := dbConn()
	delDB, err := db.Prepare("DELETE FROM Sites WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delDB.Exec(id)
	defer db.Close()
}

func monitoring(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT id, url FROM sites")
	if err != nil {
		panic(err.Error())
	}

	site := Site{}
	sites := []Site{}

	for selDB.Next() {
		var id int
		var url string
		err = selDB.Scan(&id, &url)
		if err != nil {
			panic(err.Error())
		}
		site.URL = url
		site.ID = id

		sites = append(sites, site)
	}

	for _, s := range sites {
		go ping(s, db)
	}
}

func ping(site Site, db *sql.DB) {
	time_start := time.Now()
	resp, err := http.Get(site.URL)
	if err != nil {
		insDB, err := db.Prepare("UPDATE Sites SET sucess=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insDB.Exec(false, site.ID)
		return
	}
	insDB, err := db.Prepare("UPDATE Sites SET last_execution_date=?, sucess=?, response_time=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	dt := time_start.Format(time.RFC3339)
	insDB.Exec(dt, resp.StatusCode == 200, time.Since(time_start), site.ID)

	fmt.Println(time.Since(time_start), site.URL)
}
