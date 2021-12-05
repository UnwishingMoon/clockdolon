package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/UnwishingMoon/clockdolon/pkg/app"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Start() {
	// Opening Database connection
	db, err := sql.Open("mysql", app.Conf.DB.User+":"+app.Conf.DB.Pass+"@"+"/"+app.Conf.DB.Database)
	if err != nil {
		log.Fatalf("[FATAL] Could not open connection to database: %s", err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func Close() {
	db.Close()
}
