package server

import (
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/m1cr0man/steampump/pkg/steam"
)

type GameHandler struct {
	steam *steam.API
}

func (h *GameHandler) RegisterRoutes(r *mux.Router, name string) {
	name = "/" + strings.Trim(name, "/")
	r.HandleFunc(name, h.GetGames).
		Name(name + "-list").
		Methods(http.MethodGet)
	r.HandleFunc(name+"/{appid:[0-9]+}", h.GetGame).
		Name(name).
		Methods(http.MethodGet)
	r.HandleFunc(name+"/{appid:[0-9]+}/manifest", h.GetGameManifest).
		Name(name + "-manifest").
		Methods(http.MethodGet)
	r.HandleFunc(name+"/{appid:[0-9]+}/content/{filepath:.*}", h.GetGameContent).
		Name(name + "-content").
		Methods(http.MethodGet)
}

func (h *GameHandler) GetGames(res http.ResponseWriter, req *http.Request) {
	serveJSON(res, h.steam.GetGames())
}

func (h *GameHandler) GetGame(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appID, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	serveJSON(res, h.steam.GetGame(appID))
}

func (h *GameHandler) GetGameManifest(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appID, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fullPath := h.steam.GetGameManifestPath(appID)
	serveFile(res, req, fullPath)
}

func (h *GameHandler) GetGameContent(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appID, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fullPath := path.Join(h.steam.GetGamePath(appID), vars["filepath"])
	serveFile(res, req, fullPath)
}

func NewGameHandler(steam *steam.API) *GameHandler {
	return &GameHandler{
		steam: steam,
	}
}
