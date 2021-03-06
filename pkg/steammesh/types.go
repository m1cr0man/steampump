package steammesh

import (
	"fmt"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

const DefaultPort = 9771

// Buffer size in bytes (1 MegaByte)
const IOBufferSize int = 1e+6

type Peer struct {
	Name  string                 `json:"name"`
	Proxy *httputil.ReverseProxy `json:"-"`
}

func (p *Peer) checkUrl() (*url.URL, error) {
	return url.Parse(fmt.Sprintf("http://%s:%d/", p.Name, DefaultPort))
}

func (p *Peer) Url() (url *url.URL) {
	// checkUrl will have been run when the peer was added
	url, _ = p.checkUrl()
	return url
}

type Config struct {
	Peers []Peer `json:"peers"`
}

type Action int

const (
	ActionWrite Action = iota
	ActionDelete
)

type TransferItem struct {
	Path  string      `json="path"`
	Mtime time.Time   `json="mtime"`
	Mode  os.FileMode `json="mode"`
	Size  int64       `json="size"`
}

type SyncItem struct {
	TransferItem
	Action Action `json="action"`
}

// CopierStatus Indicates the status of a GameCopier operation
type CopierStatus string

const (
	// StatusQueued Copy is queued to run when possible
	StatusQueued CopierStatus = "Queued"
	// StatusRunning Copy is running
	StatusRunning = "Running"
	// StatusCalculating Copy is running but is still counting files which must be copied
	StatusCalculating = "Calculating"
	// StatusSuccessful Copy has finished and completed successfully
	StatusSuccessful = "Successful"
	// StatusFailed Copy has finished and failed
	StatusFailed = "Failed"
)
