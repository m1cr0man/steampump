package steam

type Game struct {
	AppID              int    `acf:"appid" json:"appid"`
	Name               string `acf:"name" json:"name"`
	StateFlags         int    `acf:"StateFlags" json:"stateflags"`
	InstallDir         string `acf:"installdir" json:"installdir"`
	LastUpdated        int    `acf:"LastUpdated" json:"lastupdated"`
	UpdateResult       int    `acf:"UpdateResult" json:"updateresult"`
	SizeOnDisk         int    `acf:"SizeOnDisk" json:"sizeondisk"`
	BuildID            string `acf:"buildid" json:"buildid"`
	BytesToDownload    int    `acf:"BytesToDownload" json:"bytestodownload"`
	BytesDownloaded    int    `acf:"BytesDownloaded" json:"bytesdownloaded"`
	AutoUpdateBehavior int    `acf:"AutoUpdateBehavior" json:"autoupdatebehavior"`
}
