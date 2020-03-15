package steampump

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/m1cr0man/steampump/pkg/steam"
)

type Server struct {
	app    *App
	steam  *steam.API
	Router *mux.Router
	Port   int
}

type FileInfo struct {
	Name  string    `json="name"`
	Dir   bool      `json="dir"`
	Mtime time.Time `json="mtime"`
}

// Adapted from http.checkIfModifiedSince
// Returns true unless If-Modified-Since matches modtime
func (s *Server) checkModified(req *http.Request, modtime time.Time) bool {
	ims := req.Header.Get("If-Modified-Since")
	if ims == "" {
		return true
	}
	t, err := http.ParseTime(ims)
	if err != nil {
		return true
	}
	modtime = modtime.Truncate(time.Second)
	if modtime.Equal(t) {
		return false
	}
	return true
}

func (s *Server) sendDirJSON(res http.ResponseWriter, dirpath string) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")

	files, _ := ioutil.ReadDir(dirpath)

	fileJSON := make([]FileInfo, len(files))

	for i, file := range files {
		fileJSON[i] = FileInfo{
			Name:  file.Name(),
			Dir:   file.IsDir(),
			Mtime: file.ModTime(),
		}
	}

	err := json.NewEncoder(res).Encode(fileJSON)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) serveFile(res http.ResponseWriter, req *http.Request, fullPath string) {
	stat, err := os.Stat(fullPath)
	if os.IsNotExist(err) {
		http.NotFound(res, req)
		return
	}
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	// Override ServeFile for application/json requests to send better dir info
	if stat.IsDir() && req.Header.Get("Content-Type") == "application/json" {

		// Check modified date correctly
		mtime := stat.ModTime()
		if s.checkModified(req, mtime) == false {
			res.WriteHeader(http.StatusNotModified)
			return
		}
		res.Header().Set("Last-Modified", mtime.UTC().Format(http.TimeFormat))
		s.sendDirJSON(res, fullPath)
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

	r.HandleFunc("/steamapp", s.ServeGames).
		Name("steamapp-list").
		Methods(http.MethodGet)
	r.HandleFunc("/steamapp/{appid:[0-9]+}", s.ServeGame).
		Name("steamapp").
		Methods(http.MethodGet)
	r.HandleFunc("/steamapp/{appid:[0-9]+}/manifest", s.ServeGameManifest).
		Name("steamapp-manifest").
		Methods(http.MethodGet)
	r.HandleFunc("/steamapp/{appid:[0-9]+}/content/{filepath:.*}", s.ServeGameContent).
		Name("steamapp-content").
		Methods(http.MethodGet)

	return &s
}
