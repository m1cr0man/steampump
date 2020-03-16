package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/m1cr0man/steampump/pkg/steampump"
)

type AppHandler struct {
	app *steampump.App
}

func (h *AppHandler) RegisterRoutes(r *mux.Router, name string) {
	name = "/" + strings.Trim(name, "/")
	r.HandleFunc(name+"/config", h.GetConfig).
		Name(name + "-config").
		Methods(http.MethodGet)
	r.HandleFunc(name+"/config", h.PutConfig).
		Name(name + "-config-write").
		Methods(http.MethodPut)
}

func (h *AppHandler) GetConfig(res http.ResponseWriter, req *http.Request) {
	serveJSON(res, h.app.GetConfig())
}

func (h *AppHandler) PutConfig(res http.ResponseWriter, req *http.Request) {
	newConfig := h.app.GetConfig()
	err := json.NewDecoder(req.Body).Decode(&newConfig)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.app.SetConfig(newConfig)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusAccepted)
}

func NewAppHandler(app *steampump.App) *AppHandler {
	return &AppHandler{
		app: app,
	}
}
