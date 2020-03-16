package steammesh

import (
	"fmt"
	"net/http/httputil"
	"net/url"

	"github.com/m1cr0man/steampump/pkg/steam"
)

const DefaultPort = 9771

type API struct {
	config Config
}

func (i *API) SetConfig(config Config) (err error) {
	// Ensure lists are not nil
	if config.Peers == nil {
		config.Peers = []Peer{}
	}

	// Now all the proxies will be null. Set them up
	for _, peer := range config.Peers {
		url, err := url.Parse(fmt.Sprintf("http://%s:%d/", peer.Name, DefaultPort))
		if err != nil {
			return err
		}
		peer.Proxy = httputil.NewSingleHostReverseProxy(url)
	}

	i.config = config
	return
}

func (i *API) GetPeers() []Peer {
	return i.config.Peers
}

func (i *API) CopyGameFrom(srcHost string, srcGame, dstGame steam.Game) error {
	return nil
}

func NewAPI() *API {
	return &API{
		config: Config{
			Peers: []Peer{},
		},
	}
}
