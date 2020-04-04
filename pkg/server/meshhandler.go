package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/m1cr0man/steampump/pkg/steam"
	"github.com/m1cr0man/steampump/pkg/steammesh"
)

var client = http.Client{}

type MeshHandler struct {
	steam   *steam.API
	mesh    *steammesh.API
	copiers []*steammesh.GameCopier
}

func (h *MeshHandler) RegisterRoutes(r *mux.Router, name string) {
	name = "/" + strings.Trim(name, "/")
	r.HandleFunc(name, h.GetPeers).
		Name(name+"-list").
		Methods(http.MethodOptions, http.MethodGet)

	r.HandleFunc(name+"/copy/{peername}/{appid:[0-9]+}", h.CopyGameFrom).
		Name(name+"-copy").
		Methods(http.MethodOptions, http.MethodPost)

	r.HandleFunc(name+"/copy", h.ListCopiers).
		Name(name+"-copy-list").
		Methods(http.MethodOptions, http.MethodGet)

	// Register the Proxy Middleware. This will apply to all handlers
	r.Use(h.ProxyPeerMiddleware)
}

func (h *MeshHandler) GetPeers(res http.ResponseWriter, req *http.Request) {
	serveJSON(res, h.mesh.GetPeers())
}

func (h *MeshHandler) ListCopiers(res http.ResponseWriter, req *http.Request) {
	serveJSON(res, h.copiers)
}

func (h *MeshHandler) CopyGameFrom(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	appID, err := strconv.Atoi(vars["appid"])
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	// Find the actual peer
	peerName := vars["peername"]
	var peer steammesh.Peer
	for _, peer = range h.mesh.GetPeers() {
		if peer.Name == peerName {
			break
		}
	}

	// Fail if peer is not valid
	if peer.Name == "" {
		http.Error(res, "Invalid peer specified", http.StatusBadRequest)
		return
	}

	// Load manifest from remote system
	relURL, _ := url.Parse(fmt.Sprintf("games/%d/manifest", appID))
	mfstres, err := client.Get(peer.Url().ResolveReference(relURL).String())
	if err != nil {
		http.Error(res, "Invalid appID specified", http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(mfstres.Body)
	if err != nil {
		fmt.Printf(err.Error())
		http.Error(res, "Failed to download manifest", http.StatusInternalServerError)
		return
	}
	err = ioutil.WriteFile(h.steam.GetGameManifestPath(appID), data, 0644)
	if err != nil {
		fmt.Printf(err.Error())
		http.Error(res, "Failed to write manifest", http.StatusInternalServerError)
		return
	}
	err = h.steam.LoadGames()
	if err != nil {
		fmt.Printf(err.Error())
		http.Error(res, "Failed to reload steam games", http.StatusInternalServerError)
		return
	}

	copier := steammesh.GameCopier{
		Status:  steammesh.StatusQueued,
		AppID:   appID,
		PeerURL: *peer.Url(),
		Dest:    h.steam.GetGamePath(appID),
	}
	go copier.StartCopy()
	h.copiers = append(h.copiers, &copier)
	serveJSON(res, copier)
}

func (h *MeshHandler) ProxyPeerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		peerName := req.Header.Get("X-Peer")

		// Skip proxying if not specified
		if peerName == "" {
			next.ServeHTTP(res, req)
			return
		}

		// Find the actual peer
		var peer steammesh.Peer
		for _, peer = range h.mesh.GetPeers() {
			if peer.Name == peerName {
				break
			}
		}

		// Fail if X-Peer is specified and not valid
		if peer.Name == "" {
			http.Error(res, "Invalid peer specified", http.StatusBadRequest)
			return
		}

		// Proxy the request
		req.Header.Del("X-Peer")
		peer.Proxy.ServeHTTP(res, req)
	})
}

func NewMeshHandler(mesh *steammesh.API, steam *steam.API) *MeshHandler {
	return &MeshHandler{
		mesh:    mesh,
		steam:   steam,
		copiers: []*steammesh.GameCopier{},
	}
}
