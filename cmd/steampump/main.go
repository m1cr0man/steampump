package main

import (
	"fmt"
	"os"

	"github.com/kirsle/configdir"
	"github.com/m1cr0man/steampump/pkg/steam"
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
	app := steampump.NewApp(configDir, steam)
	err = app.LoadConfig()

	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Failed to load config: ", err)
		os.Exit(2)
		return
	}

	err = steam.SetSteamPath(app.Config.SteamPath)
	if err != nil {
		fmt.Println("Failed to set Steam path: ", err)

		if !os.IsNotExist(err) {
			os.Exit(2)
			return
		}
	} else {
		err = steam.LoadGames()
		if err != nil {
			fmt.Println("Failed to load Steam games: ", err)
			os.Exit(2)
			return
		}
	}

	server := steampump.NewServer(app, steam)
	server.Serve()
}
