package steampump

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/m1cr0man/steampump/pkg/steam"
)

const AppPort = 9771

type App struct {
	ConfigDir string
	server    *http.Server
	steam     *steam.API
}

func (a *App) configFile() string {
	return filepath.Join(a.ConfigDir, "config.json")
}

func (a *App) RunServer() {
	a.server.ListenAndServe()
}

func (a *App) StopServer() error {
	return a.server.Shutdown(context.Background())
}

func NewApp(configDir string, steam *steam.API) *App {
	return &App{
		ConfigDir: configDir,
		server:    &http.Server{Addr: fmt.Sprintf(":%d", AppPort)},
		steam:     steam,
	}
}
