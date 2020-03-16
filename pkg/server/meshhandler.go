package server

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/m1cr0man/steampump/pkg/steammesh"
)

var client = http.Client{}

type MeshHandler struct {
	mesh *steammesh.API
}

func (h *MeshHandler) RegisterRoutes(r *mux.Router, name string) {
	name = "/" + strings.Trim(name, "/")
	r.HandleFunc(name, h.GetPeers).
		Name(name + "-list").
		Methods(http.MethodGet)

	r.Use(h.ProxyPeerMiddleware)
}

func (h *MeshHandler) GetPeers(res http.ResponseWriter, req *http.Request) {
	serveJSON(res, h.mesh.GetPeers())
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
		peer.Proxy.ServeHTTP(res, req)
	})
}

func NewMeshHandler(mesh *steammesh.API) *MeshHandler {
	return &MeshHandler{
		mesh: mesh,
	}
}
