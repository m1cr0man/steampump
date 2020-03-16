package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m1cr0man/steampump/pkg/steam"
	"github.com/m1cr0man/steampump/pkg/steammesh"
	"github.com/m1cr0man/steampump/pkg/steampump"
)

const DefaultPort = 9771

type Server struct {
	app    *steampump.App
	router *mux.Router
	config Config
}

func (s *Server) SetConfig(config Config) (err error) {
	s.config = config
	return
}

func (s *Server) Serve() {
	fmt.Printf("Now listening on :%d", s.config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", s.config.Port), s.router)
}

func NewServer(app *steampump.App, mesh *steammesh.API, steam *steam.API) *Server {
	r := mux.NewRouter()
	NewAppHandler(app).RegisterRoutes(r, "app")
	NewMeshHandler(mesh).RegisterRoutes(r, "mesh")
	NewGameHandler(steam).RegisterRoutes(r, "games")

	s := Server{
		app:    app,
		router: r,
		config: Config{
			Port: DefaultPort,
		},
	}
	return &s
}
