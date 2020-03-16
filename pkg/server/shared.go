package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func serveJSON(res http.ResponseWriter, data interface{}) {
	res.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(res).Encode(data)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

// Adapted from http.checkIfModifiedSince
// Returns true unless If-Modified-Since matches modtime
func checkModified(req *http.Request, modtime time.Time) bool {
	ims := req.Header.Get("If-Modified-Since")
	if ims == "" {
		return true
	}
	t, err := http.ParseTime(ims)
	if err != nil {
		return true
	}
	modtime = modtime.Truncate(time.Second)
	return !modtime.Equal(t)
}

func sendDirJSON(res http.ResponseWriter, dirpath string) {
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

func serveFile(res http.ResponseWriter, req *http.Request, fullPath string) {
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
		if checkModified(req, mtime) == false {
			res.WriteHeader(http.StatusNotModified)
			return
		}
		res.Header().Set("Last-Modified", mtime.UTC().Format(http.TimeFormat))
		sendDirJSON(res, fullPath)
		return
	}
	http.ServeFile(res, req, fullPath)
}
