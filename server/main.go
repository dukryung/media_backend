package main

import (
	"github.com/dukryung/media_backend/server/app"
	"os"
	"os/signal"
	"syscall"
)

const configPath = "./config.json"

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	app := app.NewApp(configPath)

	err := app.RunServers()
	if err != nil {
		panic(err)
	}

	<- quit

	app.CloseServers()

}
