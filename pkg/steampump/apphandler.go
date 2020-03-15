package steampump

import (
	"encoding/json"
	"net/http"
)

func (s *Server) ServeConfig(res http.ResponseWriter, req *http.Request) {
	jsonData, err := json.Marshal(s.app.Config)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func (s *Server) WriteConfig(res http.ResponseWriter, req *http.Request) {
	var newConfig Config
	err := json.NewDecoder(req.Body).Decode(&newConfig)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Try setting the Steam path before saving the config since it does some error checking
	err = s.steam.SetSteamPath(newConfig.SteamPath)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.app.SaveConfig(newConfig)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}
