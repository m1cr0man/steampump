package steam

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/oleiade/reflections.v1"
)

type API struct {
	steamPath string
	games     []Game
}

func (i *API) SetSteamPath(steamPath string) (err error) {
	_, err = os.Stat(steamPath)

	if err != nil {
		return
	}

	i.steamPath = steamPath
	return
}

func (i *API) SteamAppsPath() string {
	return path.Join(i.steamPath, "steamapps")
}

func (i *API) GetGame(appid int) Game {
	for _, game := range i.games {
		if game.AppID == appid {
			return game
		}
	}
	return Game{}
}

func (i *API) GetGamePath(appid int) string {
	game := i.GetGame(appid)
	return path.Join(i.SteamAppsPath(), "common", game.InstallDir)
}

func (i *API) GetGameManifestPath(appid int) string {
	return path.Join(i.SteamAppsPath(), fmt.Sprintf("appmanifest_%d.acf", appid))
}

func (i *API) LoadManifest(filename string) (game Game, err error) {
	// Get maps of ACF keys to fields, and fields to types
	fields, _ := reflections.Fields(game)
	typesMap := make(map[string]reflect.Kind, len(fields))
	tagsMap := make(map[string]string, len(fields))
	var field, tag string
	for _, field = range fields {
		tag, _ = reflections.GetFieldTag(game, field, "acf")
		tagsMap[tag] = field
		typesMap[field], _ = reflections.GetFieldKind(game, field)
	}

	// Read the acf, parse the lines out
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	var key string
	var char byte
	var builder strings.Builder
	var inQuotes bool = false
	for _, char = range data {
		switch char {
		case '"':
			// Look for quotes
			inQuotes = !inQuotes

			if !inQuotes {
				value := builder.String()

				// First pair of quotes contains the key
				// Map the acf tag to the actual struct field name
				if fieldName, ok := tagsMap[value]; ok {
					key = fieldName

					// The second is the value
				} else if kind, ok := typesMap[key]; ok {
					// Could be int or string - make the type conversion
					if kind == reflect.Int {
						var intval int64
						intval, err = strconv.ParseInt(value, 10, 0)
						if err != nil {
							return
						}
						reflections.SetField(&game, key, int(intval))
					} else {
						reflections.SetField(&game, key, value)
					}
				}
			}
			builder.Reset()

		case '\n':
			// Reset on new lines
			inQuotes = false
			key = ""
			builder.Reset()

		default:
			// Save regular characters
			builder.WriteByte(char)
		}
	}
	return
}

func (i *API) LoadGames() (err error) {
	var games []Game
	var game Game
	files, err := ioutil.ReadDir(i.SteamAppsPath())
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".acf") && !file.IsDir() {
			game, err = i.LoadManifest(path.Join(i.SteamAppsPath(), file.Name()))
			if err != nil {
				return
			}
			games = append(games, game)
		}
	}

	i.games = games
	return
}

func NewAPI() *API {
	return &API{
		games: []Game{},
	}
}
