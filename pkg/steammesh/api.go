package steammesh

import (
	"net"
	"net/http"
	"net/http/httputil"
	"time"
)

type API struct {
	config Config
}

func (i *API) SetConfig(config Config) (err error) {
	// Ensure lists are not nil
	if config.Peers == nil {
		config.Peers = []Peer{}
	}

	// Now all the proxies will be null. Set them up
	for i, peer := range config.Peers {
		url, err := peer.checkUrl()
		if err != nil {
			return err
		}
		peer.Proxy = httputil.NewSingleHostReverseProxy(url)
		peer.Proxy.Transport = &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: time.Second,
			}).DialContext,
		}
		config.Peers[i] = peer
	}

	i.config = config
	return
}

func (i *API) GetPeers() []Peer {
	return i.config.Peers
}

func NewAPI() *API {
	return &API{
		config: Config{
			Peers: []Peer{},
		},
	}
}
