package steampump

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/m1cr0man/steampump/pkg/steam"
)

type Server struct {
	steam *steam.API
}

func (s *Server) serveFile(res http.ResponseWriter, req *http.Request, fullPath string) {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(res, req)
		return
	}
	http.ServeFile(res, req, fullPath)
}

func NewServer(steam *steam.API) {
	s := Server{steam: steam}
	r := mux.NewRouter()

	r.HandleFunc("/steamapp/{appid:[0-9]+}/?", s.ServeGame).
		Name("steamapp").
		Methods(http.MethodGet)
	r.HandleFunc("/steamapp/{appid:[0-9]+}/manifest/?", s.ServeGameManifest).
		Name("steamapp-manifest").
		Methods(http.MethodGet)
	r.HandleFunc("/steamapp/{appid:[0-9]+}/content/{filepath:.*}", s.ServeGameContent).
		Name("steamapp-content").
		Methods(http.MethodGet)
}
