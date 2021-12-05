package app

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var Conf = appConf{}

type appConf struct {
	Bot botConf      `toml:"bot"`
	DB  databaseConf `toml:"database"`
}

type botConf struct {
	Prefix string   `toml:"prefix"`
	Token  string   `toml:"token"`
	Admins []string `toml:"admins"`
}

type databaseConf struct {
	Database string `toml:"database"`
	User     string `toml:"user"`
	Pass     string `toml:"pass"`
}

func init() {
	dir, err := os.Getwd()
	if err == nil {
		dir = dir + "/"
	}

	if _, err := toml.DecodeFile(dir+"config.toml", &Conf); err != nil {
		log.Fatalf("[FATAL] Could not read %s file: %s", "config.toml", err.Error())
	}
}
