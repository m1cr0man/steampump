package steammesh

import "net/http/httputil"

type Peer struct {
	Name  string `json:"name"`
	Proxy *httputil.ReverseProxy
}

type Config struct {
	Peers []Peer `json:"peers"`
}
