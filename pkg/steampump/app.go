package steampump

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/m1cr0man/steampump/pkg/steam"
)

const AppPort = 9771

type App struct {
	ConfigDir string
	Config    *Config
	steam     *steam.API
}

func (a *App) configFile() string {
	return filepath.Join(a.ConfigDir, "config.json")
}

func (a *App) LoadConfig() (err error) {
	var data []byte
	data, err = ioutil.ReadFile(a.configFile())
	if err != nil {
		return
	}

	err = json.Unmarshal(data, a.Config)
	return
}

func (a *App) SaveConfig(config Config) (err error) {
	var data []byte
	data, err = json.Marshal(config)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(a.configFile(), data, 0644)
	if err != nil {
		return
	}
	a.Config = &config
	return
}

func NewApp(configDir string, steam *steam.API) *App {
	return &App{
		ConfigDir: configDir,
		Config:    &Config{},
		steam:     steam,
	}
}
