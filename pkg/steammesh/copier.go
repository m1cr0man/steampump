package steammesh

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"
)

type GameCopier struct {
	Running bool    `json:"running"`
	AppID   int     `json:"appid"`
	PeerURL url.URL `json:"peer"`
	Dest    string  `json:"dest"`
}

func (g *GameCopier) getPaths(remoteURL *url.URL, item *SyncItem) (string, *url.URL) {
	localPath := path.Join(g.Dest, item.Path)
	relPath, _ := url.Parse(item.Path)
	remotePath := remoteURL.ResolveReference(relPath)
	return localPath, remotePath
}

func (g *GameCopier) DiffDirectory(remoteURL *url.URL, item *SyncItem) (items []*SyncItem, err error) {
	localPath, remotePath := g.getPaths(remoteURL, item)
	req, _ := http.NewRequest(http.MethodGet, remotePath.String(), nil)
	req.Header.Add("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	remoteFiles := []TransferItem{}
	err = json.NewDecoder(res.Body).Decode(&remoteFiles)
	if err != nil {
		return
	}

	localFiles, err := ioutil.ReadDir(localPath)
	if os.IsNotExist(err) {
		err = nil
		localFiles = []os.FileInfo{}
	}
	if err != nil {
		return
	}

	// File paths from remote are just file names relative to the requested folder
	remoteMap := make(map[string]TransferItem, len(remoteFiles))
	for _, file := range remoteFiles {
		remoteMap[file.Path] = file
	}

	localMap := make(map[string]os.FileInfo, len(localFiles))
	for _, file := range localFiles {
		localMap[file.Name()] = file
	}

	for _, file := range localFiles {
		// Delete
		if _, ok := remoteMap[file.Name()]; !ok {
			localFile := path.Join(item.Path, file.Name()) // TODO check
			items = append(items, &SyncItem{
				TransferItem: TransferItem{
					Path:  localFile,
					Mtime: file.ModTime(),
					Mode:  file.Mode(),
					Size:  file.Size(),
				},
				Action: ActionDelete,
			})
		}
	}
	for _, ritem := range remoteFiles {
		// Create or modify
		if file, ok := localMap[ritem.Path]; !(ok && file.ModTime().Equal(ritem.Mtime)) {
			// Make ritem.Path relative to requested directory
			// It is just the file name at the moment
			ritem.Path = path.Join(item.Path, ritem.Path)
			items = append(items, &SyncItem{
				TransferItem: ritem,
				Action:       ActionWrite,
			})
		}
	}

	return
}

func (g *GameCopier) DownloadItem(remoteURL *url.URL, item *SyncItem) (err error) {
	localPath, remotePath := g.getPaths(remoteURL, item)
	fd, err := os.Create(localPath)
	if err != nil {
		return
	}
	defer fd.Close()

	res, err := http.Get(remotePath.String())
	if err != nil {
		return
	}
	defer res.Body.Close()

	buf := make([]byte, IOBufferSize)
	var readerr error
	var n int = IOBufferSize
	for n == IOBufferSize {
		n, readerr = res.Body.Read(buf)
		// Skip checking readerr, it might just be EOF

		// Write returns an error if num bytes written != len of input
		_, err = fd.Write(buf[:n])
		if err != nil {
			return
		}

		if readerr != nil && readerr != io.EOF {
			return readerr
		}
	}

	err = os.Chmod(localPath, item.Mode.Perm())
	return
}

func (g *GameCopier) CopyGameFrom() (err error) {
	relPath, _ := url.Parse(fmt.Sprintf("games/%d/content/", g.AppID))
	remoteURL := g.PeerURL.ResolveReference(relPath)

	// Create an initial sync item for the root of the game
	items := []*SyncItem{
		&SyncItem{
			TransferItem: TransferItem{
				Path:  ".",
				Mtime: time.Now(),
				Mode:  0755 | os.ModeDir,
				Size:  0,
			},
			Action: ActionWrite,
		},
	}
	var newItems []*SyncItem
	var item *SyncItem

	for len(items) > 0 {
		// Pop sync item
		item = items[0]
		items = items[1:]
		localPath := path.Join(g.Dest, item.Path)
		fmt.Println(item.Path)

		// Write or delete? write covers syncing folders
		switch item.Action {
		case ActionWrite:
			// Directory? create it and sync
			if item.Mode.IsDir() {
				if err = os.MkdirAll(localPath, item.Mode); err != nil {
					return
				}
				newItems, err = g.DiffDirectory(remoteURL, item)
				items = append(items, newItems...)

				// File? Download it
			} else {
				err = g.DownloadItem(remoteURL, item)
			}

			if err != nil {
				return
			}
			err = os.Chtimes(localPath, item.Mtime, item.Mtime)

		case ActionDelete:
			err = os.RemoveAll(localPath)
		}

		if err != nil {
			return
		}
	}

	return
}

func (g *GameCopier) StartCopy() {
	fmt.Printf("Starting copy: %d from %s to %s\n", g.AppID, g.PeerURL.String(), g.Dest)
	g.Running = true
	err := g.CopyGameFrom()
	g.Running = false

	if err != nil {
		fmt.Println("Copy ", g.AppID, " failed with error: ", err)
	} else {
		fmt.Println("Copy ", g.AppID, " completed successfully")
	}
}