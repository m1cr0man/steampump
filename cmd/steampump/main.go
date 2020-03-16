package main

import (
	"fmt"
	"os"

	"github.com/kirsle/configdir"
	"github.com/m1cr0man/steampump/pkg/server"
	"github.com/m1cr0man/steampump/pkg/steam"
	"github.com/m1cr0man/steampump/pkg/steammesh"
	"github.com/m1cr0man/steampump/pkg/steampump"
)

func main() {
	var err error

	configDir := configdir.LocalConfig("steampump")
	err = configdir.MakePath(configDir) // Ensure it exists.
	if err != nil {
		fmt.Println("Failed to get config dir: ", err)
		os.Exit(1)
		return
	}

	steam := steam.NewAPI()
	mesh := steammesh.NewAPI()
	app := steampump.NewApp(configDir, steam, mesh)

	if err = app.LoadConfig(); err != nil && !os.IsNotExist(err) {
		fmt.Println("Failed to load config: ", err)
		os.Exit(2)
		return
	}

	if err = steam.LoadGames(); err != nil && app.GetConfig().Steam.SteamPath != "" {
		fmt.Println("Failed to load Steam games: ", err)
		os.Exit(2)
		return
	}

	server := server.NewServer(app, mesh, steam)
	server.Serve()
}
