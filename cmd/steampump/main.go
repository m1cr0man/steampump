package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/kirsle/configdir"
	"github.com/m1cr0man/steampump/pkg/server"
	"github.com/m1cr0man/steampump/pkg/steam"
	"github.com/m1cr0man/steampump/pkg/steammesh"
	"github.com/m1cr0man/steampump/pkg/steampump"
	"github.com/pkg/browser"
)

func main() {
	systray.RunWithAppWindow("Steampump", 1024, 768, onReady, onExit)
	// systray.Run(onReady, onExit)

}

func onReady() {
	IconParsed, _ := base64.StdEncoding.DecodeString(Icon)
	systray.SetIcon(IconParsed)
	systray.SetTitle("Steampump")
	systray.SetTooltip("Steampump")
	mUi := systray.AddMenuItem("Open", "Open the UI")
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				fmt.Println("Requesting quit")
				systray.Quit()
				fmt.Println("Finished quitting")
				os.Exit(0)
			case <-mUi.ClickedCh:
				err := browser.OpenURL("http://localhost:9771/ui/")
				if err != nil {
					fmt.Printf(err.Error())
				}
			}
		}
	}()

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
	go server.Serve()
}

func onExit() {
	systray.Quit()
}
