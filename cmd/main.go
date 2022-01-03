package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/UnwishingMoon/clockdolon/pkg/bot"
	"github.com/UnwishingMoon/clockdolon/pkg/cetus"
	"github.com/UnwishingMoon/clockdolon/pkg/db"
)

func main() {
	// Starts cetus timers and populate structs
	cetus.Start()

	// Opens Database connection
	db.Start()

	// Starting discord bot
	bot.Start()
	defer bot.Close()

	// Waiting for a signal to exit
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
