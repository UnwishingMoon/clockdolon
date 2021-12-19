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
	var err error
	db, err = sql.Open("mysql", app.Conf.DB.User+":"+app.Conf.DB.Pass+"@/"+app.Conf.DB.Database)

	if err != nil {
		log.Fatalf("[FATAL] Could not open connection to database: %s", err.Error())
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("[FATAL] Could not open connect to database: %s", err.Error())
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
}

func UserAlertExist(guild string, user string) bool {
	row, err := db.Query("SELECT AL_COD FROM alerts WHERE AL_GUILD=? AND AL_USER=? AND AL_DISABLED=0 LIMIT 1", guild, user)
	if err != nil {
		panic(err)

	}
	defer row.Close()

	if row.Next() {
		return true
	}

	return false
}

func AddUserAlert(guild string, user string, minutes int) error {
	_, err := db.Exec("INSERT INTO alerts SET AL_GUILD=?, AL_USER=?, AL_TIME=?", guild, user, minutes)
	return err
}

func RemoveUserAlert(guild string, user string) error {
	_, err := db.Exec("UPDATE alerts SET AL_DISABLED=1 WHERE AL_GUILD=? AND AL_USER=?", guild, user)
	return err
}
