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

func UserAlertExist(guild string, user string) bool {
	var e int
	if err := db.QueryRow("SELECT AL_COD FROM alerts WHERE AL_GUILD=? AND AL_USER=?", guild, user).Scan(&e); err != nil {
		if err == sql.ErrNoRows {
			return true
		}
		return false
	}

	return false
}

func AddUserAlert(guild string, user string, minutes int) error {
	_, err := db.Exec("INSERT INTO alerts SET AL_GUILD=?, AL_USER=?, AL_TIME=?", guild, user, minutes)
	return err
}
