package steampump

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/m1cr0man/steampump/pkg/steam"
	"github.com/m1cr0man/steampump/pkg/steammesh"
)

type App struct {
	configDir string
	config    Config
	steam     *steam.API
	mesh      *steammesh.API
}

func (a *App) configFile() string {
	return filepath.Join(a.configDir, "config.json")
}

func (a *App) LoadConfig() (err error) {
	var data []byte
	data, err = ioutil.ReadFile(a.configFile())
	if err != nil {
		return
	}

	newConfig := a.config
	if err = json.Unmarshal(data, &newConfig); err != nil {
		return
	}

	err = a.SetConfig(newConfig)
	return
}

func (a *App) SetConfig(config Config) (err error) {

	if err = a.steam.SetConfig(config.Steam); err != nil {
		return
	}

	if err = a.mesh.SetConfig(config.Mesh); err != nil {
		// Revert steam config
		a.steam.SetConfig(a.config.Steam)
		return
	}

	var data []byte
	if data, err = json.Marshal(config); err != nil {
		// Revert others
		a.mesh.SetConfig(a.config.Mesh)
		a.steam.SetConfig(a.config.Steam)
		return
	}

	if err = ioutil.WriteFile(a.configFile(), data, 0644); err != nil {
		// Revert others
		a.mesh.SetConfig(a.config.Mesh)
		a.steam.SetConfig(a.config.Steam)
		return
	}

	a.config = config
	return
}

func (a *App) GetConfig() Config {
	return a.config
}

func NewApp(configDir string, steamapi *steam.API, mesh *steammesh.API) *App {
	return &App{
		configDir: configDir,
		config:    Config{},
		steam:     steamapi,
		mesh:      mesh,
	}
}
