package app

import (
	"github.com/dukryung/media_backend/server/media"
	"github.com/dukryung/media_backend/server/types"
)

type App struct {
	servers []types.Server
	appConfig types.AppConfig
	mediaServer *media.Server
}

func NewApp(configPath string) *App {
	app := App{}

	app.appConfig = types.AppConfig{}
	err := app.appConfig.LoadAppConfig(configPath)
	if err != nil {
		panic(err)
	}

	app.initServers()

	return &app
}

func (app *App) initServers() {
	app.mediaServer = media.NewServer(app.appConfig)

	media.RegisterMediaServer(app.mediaServer.GrpcServer,app.mediaServer)
	app.servers = append(app.servers, app.mediaServer)

}

func (app *App) RunServers() error {
	for _, server := range app.servers {
		go server.Run()
	}
	return nil
}

func (app *App) CloseServers() {
	for _, server := range app.servers {
		server.Close()
	}
}
