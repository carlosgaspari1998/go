package database

import (
	"fmt"
	"net/http"
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
	selDB, err := db.db.Query("SELECT id, url FROM sites")
	if err != nil {
		panic(err.Error())
	}

	site := models.Site{}
	sites := []models.Site{}

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

func ping(site models.Site, db *DB) {
	time_start := time.Now()
	resp, err := http.Get(site.URL)
	if err != nil {
		insDB, err := db.db.Prepare("UPDATE Sites SET sucess=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}
		insDB.Exec(false, site.ID)
		return
	}
	insDB, err := db.db.Prepare("UPDATE Sites SET last_execution_date=?, sucess=?, response_time=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	dt := time_start.Format(time.RFC3339)
	insDB.Exec(dt, resp.StatusCode == 200, time.Since(time_start), site.ID)
	fmt.Println(time.Since(time_start), site.URL)
}
