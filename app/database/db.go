package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func Open() *sql.DB {
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
