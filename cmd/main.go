package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/UnwishingMoon/clockdolon/pkg/bot"
	"github.com/UnwishingMoon/clockdolon/pkg/cetus"
)

func main() {
	cetus.PopulateCetusTime()

	// Opens Database connection
	//db.Start()
	//defer db.Close()

	// Starting discord bot
	dg, err := bot.Start()
	if err != nil {
		log.Fatalf("[FATAL] Error during bot initialization: %s", err.Error())
	}
	defer dg.Close()

	// Waiting for a signal to exit
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
