package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"siteapi/app/database"
	"time"
)

type Site struct {
	ID                  int        `json:"id"`
	URL                 string     `json:"url"`
	Frequency           float64    `json:"frequency"`
	LastExecutionDate   *time.Time `json:"last_execution_date"`
	Sucess              *bool      `json:"sucess"`
	ResponseTime        *float64   `json:"response_time"`
	ResponseAverageTime *float64   `json:"response_average_time"`
	CreationDate        time.Time  `json:"creation_date"`
	LastUpdatedDate     *time.Time `json:"last_updated_date"`
}

func GetSites() ([]Site, error) {
	db := database.Open()
	selDB, err := db.Query("SELECT id, url, frequency, sucess, response_time, response_average_time, creation_date, last_execution_date, last_updated_date FROM sites")
	if err != nil {
		panic(err.Error())
	}
	site := Site{}
	sites := []Site{}
	for selDB.Next() {
		err = selDB.Scan(&site.ID, &site.URL, &site.Frequency, &site.Sucess, &site.ResponseTime, &site.ResponseAverageTime, &site.CreationDate, &site.LastExecutionDate, &site.LastUpdatedDate)
		if err != nil {
			panic(err.Error())
		}
		sites = append(sites, site)
	}
	return sites, nil
}

func GetSiteById(id string) (Site, error) {
	db := database.Open()
	selDB := db.QueryRow("SELECT id, url, frequency, last_execution_date, sucess, response_time, response_average_time, creation_date, last_updated_date FROM sites WHERE id=?", id)
	site := Site{}

	err := selDB.Scan(&site.ID, &site.URL, &site.Frequency, &site.LastExecutionDate, &site.Sucess, &site.ResponseTime, &site.ResponseAverageTime, &site.CreationDate, &site.LastUpdatedDate)
	if err != nil {
		panic(err.Error())
	}
	return site, nil
}

func CreateSite(site Site) error {
	db := database.Open()
	insDB, err := db.Prepare("INSERT INTO sites(url, frequency, creation_date) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insDB.Exec(site.URL, site.Frequency, site.CreationDate)
	return err
}

func UpdateSite(site Site) error {
	db := database.Open()
	time_start := time.Now()
	updDB, err := db.Prepare("UPDATE Sites SET url=?, frequency=?, creation_date=?, last_updated_date=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	updDB.Exec(site.URL, site.Frequency, site.CreationDate, time_start.Format(time.RFC3339), site.ID)
	return err
}

func DeleteSite(id string) error {
	db := database.Open()
	delDB, err := db.Prepare("DELETE FROM Sites WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delDB.Exec(id)
	return err
}

func MonitoringSite() {
	db := database.Open()
	sites := make(chan Site, 5)
	readSites := make(chan Site, 5)

	go GetSitesMonitoring(sites, readSites, db)
	go Ping(sites, readSites, db)
}

func GetSitesMonitoring(sites chan<- Site, readSites <-chan Site, db *sql.DB) {
	selDB, err := db.Query("SELECT id, url, frequency FROM sites")
	if err != nil {
		panic(err.Error())
	}

	site := Site{}
	for selDB.Next() {
		err = selDB.Scan(&site.ID, &site.URL, &site.Frequency)
		if err != nil {
			panic(err.Error())
		}
		sites <- site
	}

	for rs := range readSites {
		sites <- rs
	}
}

func Ping(sites <-chan Site, readSites chan<- Site, db *sql.DB) {
	for site := range sites {
		time_start := time.Now()
		resp, errRequest := http.Get(site.URL)
		updDB, err := db.Prepare("UPDATE Sites SET last_execution_date=?, sucess=?, response_time=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		lastExecutionDate := time_start.Format(time.RFC3339)
		sucess := (errRequest == nil && resp.StatusCode == 200)

		updDB.Exec(lastExecutionDate, sucess, time.Since(time_start), site.ID)

		if !sucess {
			go SaveLogErrorInFile(site.URL)
		}
		go AddSiteMonitoring(site, readSites)
	}
}

func SaveLogErrorInFile(url string) {
	f, err := os.OpenFile("errorRequest", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	logger := log.New(f, "Error - ", log.LstdFlags)
	logger.Println("- URL:", url)
}

func AddSiteMonitoring(site Site, readSites chan<- Site) {
	fmt.Println(site.URL)
	time.Sleep(time.Millisecond * time.Duration(site.Frequency))
	readSites <- site
}
