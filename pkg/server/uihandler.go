package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/markbates/pkger"
)

type UIHandler struct {
}

func (h *UIHandler) RegisterRoutes(r *mux.Router, name string) {
	name = "/" + strings.Trim(name, "/")
	r.PathPrefix(name).Handler(http.StripPrefix(name, http.FileServer(pkger.Dir("/ui/build"))))
}

func NewUIHandler() *UIHandler {
	return &UIHandler{}
}
