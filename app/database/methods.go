package database

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"siteapi/app/models"
	"time"
)

func (db *DB) GetSites() ([]models.Site, error) {
	selDB, err := db.db.Query("SELECT id, url, frequency, sucess, response_time, response_average_time, creation_date, last_execution_date, last_updated_date FROM sites")
	if err != nil {
		panic(err.Error())
	}
	site := models.Site{}
	sites := []models.Site{}
	for selDB.Next() {
		err = selDB.Scan(&site.ID, &site.URL, &site.Frequency, &site.Sucess, &site.ResponseTime, &site.ResponseAverageTime, &site.CreationDate, &site.LastExecutionDate, &site.LastUpdatedDate)
		if err != nil {
			panic(err.Error())
		}
		sites = append(sites, site)
	}
	return sites, nil
}

func (db *DB) GetSiteById(id string) (models.Site, error) {
	selDB := db.db.QueryRow("SELECT id, url, frequency, last_execution_date, sucess, response_time, response_average_time, creation_date, last_updated_date FROM sites WHERE id=?", id)
	site := models.Site{}

	err := selDB.Scan(&site.ID, &site.URL, &site.Frequency, &site.LastExecutionDate, &site.Sucess, &site.ResponseTime, &site.ResponseAverageTime, &site.CreationDate, &site.LastUpdatedDate)
	if err != nil {
		panic(err.Error())
	}
	return site, nil
}

func (db *DB) CreateSite(site models.Site) error {
	insDB, err := db.db.Prepare("INSERT INTO sites(url, frequency, creation_date) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insDB.Exec(site.URL, site.Frequency, site.CreationDate)
	return err
}

func (db *DB) UpdateSite(site models.Site) error {
	time_start := time.Now()
	updDB, err := db.db.Prepare("UPDATE Sites SET url=?, frequency=?, creation_date=?, last_updated_date=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	updDB.Exec(site.URL, site.Frequency, site.CreationDate, time_start.Format(time.RFC3339), site.ID)
	return err
}

func (db *DB) DeleteSite(id string) error {
	delDB, err := db.db.Prepare("DELETE FROM Sites WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delDB.Exec(id)
	return err
}

func (db *DB) MonitoringSite() {
	sites := make(chan models.Site, 5)
	readSites := make(chan models.Site, 5)

	go GetSitesMonitoring(sites, readSites, db)
	go Ping(sites, readSites, db)
}

func GetSitesMonitoring(sites chan<- models.Site, readSites <-chan models.Site, db *DB) {
	selDB, err := db.db.Query("SELECT id, url, frequency FROM sites")
	if err != nil {
		panic(err.Error())
	}

	site := models.Site{}
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

func Ping(sites <-chan models.Site, readSites chan<- models.Site, db *DB) {
	for site := range sites {
		time_start := time.Now()
		resp, errRequest := http.Get(site.URL)
		updDB, err := db.db.Prepare("UPDATE Sites SET last_execution_date=?, sucess=?, response_time=? WHERE id=?")
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

func AddSiteMonitoring(site models.Site, readSites chan<- models.Site) {
	fmt.Println(site.URL)
	time.Sleep(time.Millisecond * time.Duration(site.Frequency))
	readSites <- site
}
