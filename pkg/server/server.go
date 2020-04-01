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

func SetCorsHeaders(res http.ResponseWriter) {
	res.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	res.Header().Add("Access-Control-Allow-Headers", "Accept")
	res.Header().Add("Access-Control-Allow-Headers", "X-Peer")
	// TODO dev mode option
	res.Header().Set("Access-Control-Allow-Origin", fmt.Sprintf("http://localhost:%d", 3000))
	// Needed for caching
	res.Header().Set("Vary", "Origin")
}

func CorsHandler(res http.ResponseWriter, req *http.Request) {
	SetCorsHeaders(res)
	res.WriteHeader(http.StatusNoContent)
	return
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		SetCorsHeaders(res)
		if req.Method == http.MethodOptions {
			res.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(res, req)
	})
}

func NewServer(app *steampump.App, mesh *steammesh.API, steam *steam.API) *Server {
	r := mux.NewRouter()
	NewAppHandler(app).RegisterRoutes(r, "app")
	NewMeshHandler(mesh, steam).RegisterRoutes(r, "mesh")
	NewGameHandler(steam).RegisterRoutes(r, "games")
	r.Use(CorsMiddleware)

	s := Server{
		app:    app,
		router: r,
		config: Config{
			Port: DefaultPort,
		},
	}
	return &s
}
