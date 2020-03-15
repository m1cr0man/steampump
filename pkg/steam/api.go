package steam

import (
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

func (i *API) steamAppsPath() string {
	return path.Join(i.steamPath, "steamapps")
}

func (i *API) LoadManifest(filename string) (game Game, err error) {
	// Get maps of ACF keys to fields, and fields to types
	fields, _ := reflections.Fields(game)
	fieldMap := make(map[string]reflect.Kind, len(fields))
	tagsMap := make(map[string]string, len(fields))
	var field, tag string
	for _, field = range fields {
		tag, _ = reflections.GetFieldTag(game, field, "acf")
		tagsMap[tag] = field
		fieldMap[field], _ = reflections.GetFieldKind(game, field)
	}

	// Read the acf, parse the lines out
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	var key string
	var char byte
	var builder strings.Builder
	var group int
	inQuotes := false
	for _, char = range data {
		switch char {
		case '"':
			// Look for quotes
			inQuotes = !inQuotes

			if !inQuotes {
				value := builder.String()

				// First pair of quotes contains the key
				// Map the acf tag to the actual struct field name
				if fieldName, ok := tagsMap[value]; group == 0 && ok {
					key = fieldName
					group++

					// The second is the value
				} else if kind, ok := fieldMap[key]; group == 1 && ok {
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
					group++
				}
			}
			builder.Reset()

		case '\n':
			// Reset on new lines
			inQuotes = false
			group = 0
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
	files, err := ioutil.ReadDir(i.steamAppsPath())
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.Contains(file.Name(), ".acf") && !file.IsDir() {
			game, err = i.LoadManifest(path.Join(i.steamAppsPath(), file.Name()))
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
