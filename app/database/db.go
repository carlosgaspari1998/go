package database

import (
	"database/sql"
	"log"
	"siteapi/app/models"

	_ "github.com/go-sql-driver/mysql"
)

type SiteDB interface {
	Open() error
	Close() error
	GetSites() ([]models.Site, error)
	GetSiteById(string) (models.Site, error)
	CreateSite(models.Site) error
	UpdateSite(models.Site) error
	DeleteSite(string) error
	MonitoringSite()
}

type DB struct {
	db *sql.DB
}

func (d *DB) Open() error {
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName+"?parseTime=true")
	if err != nil {
		return err
	}
	log.Println("Connected to Database!")
	d.db = db
	return nil
}

func (d *DB) Close() error {
	return d.db.Close()
}
