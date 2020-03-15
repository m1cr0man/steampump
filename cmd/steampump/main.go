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
	steam.SetSteamPath("D:\\program files (x86)\\steam")
	steam.LoadGames()
	app := steampump.NewApp(configDir, steam)
	app.RunServer()
}
