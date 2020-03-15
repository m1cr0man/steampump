package steampump

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func (s *Server) ServeGame(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	game := s.steam.GetGame(vars["appid"])

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
	fullPath := s.steam.GetGameManifestPath(vars["appid"])
	s.serveFile(res, req, fullPath)
}

func (s *Server) ServeGameContent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fullPath := path.Join(s.steam.GetGamePath(vars["appid"]), vars["filepath"])
	s.serveFile(res, req, fullPath)
}
