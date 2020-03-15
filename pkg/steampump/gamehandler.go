package steampump

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"

	"github.com/gorilla/mux"
)

func (s *Server) ServeGame(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appID, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	game := s.steam.GetGame(appID)
	jsonData, err := json.Marshal(game)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(jsonData)
}

func (s *Server) ServeGameManifest(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appID, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fullPath := s.steam.GetGameManifestPath(appID)
	s.serveFile(res, req, fullPath)
}

func (s *Server) ServeGameContent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appID, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fullPath := path.Join(s.steam.GetGamePath(appID), vars["filepath"])
	s.serveFile(res, req, fullPath)
}
