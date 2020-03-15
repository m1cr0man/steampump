package steampump

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/m1cr0man/steampump/pkg/steam"
)

type Server struct {
	app    *App
	steam  *steam.API
	Router *mux.Router
	Port   int
}

func (s *Server) serveFile(res http.ResponseWriter, req *http.Request, fullPath string) {
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		http.NotFound(res, req)
		return
	}
	http.ServeFile(res, req, fullPath)
}

func (s *Server) Serve() {
	fmt.Printf("Now listening on :%d", s.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Router)
}

func NewServer(app *App, steam *steam.API) *Server {
	r := mux.NewRouter()
	s := Server{app: app, steam: steam, Router: r, Port: 9771}

	r.HandleFunc("/app/config", s.ServeConfig).
		Name("appconfig").
		Methods(http.MethodGet)
	r.HandleFunc("/app/config", s.WriteConfig).
		Name("appconfig-write").
		Methods(http.MethodPut)

	r.HandleFunc("/steamapp/{appid:[0-9]+}/?", s.ServeGame).
		Name("steamapp").
		Methods(http.MethodGet)
	r.HandleFunc("/steamapp/{appid:[0-9]+}/manifest/?", s.ServeGameManifest).
		Name("steamapp-manifest").
		Methods(http.MethodGet)
	r.HandleFunc("/steamapp/{appid:[0-9]+}/content/{filepath:.*}", s.ServeGameContent).
		Name("steamapp-content").
		Methods(http.MethodGet)

	return &s
}
