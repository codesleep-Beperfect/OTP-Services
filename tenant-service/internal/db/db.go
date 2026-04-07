package db

import (
	"database/sql"
	"log"
	"tenant-service/config"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func Init() {
	var err error
	DB, err = sql.Open("mysql", config.App.MySQL.DSN)
	if err != nil {
		log.Fatal(err)
	}
}