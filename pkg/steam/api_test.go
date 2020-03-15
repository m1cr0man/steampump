package steam

import (
	"reflect"
	"testing"
)

func TestAPI_LoadManifest(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name     string
		i        *API
		args     args
		wantGame Game
		wantErr  bool
	}{
		{
			name: "Can read app manifest",
			i:    &API{},
			args: args{
				filename: "D:\\program files (x86)\\steam\\steamapps\\appmanifest_218.acf",
			},
			wantGame: Game{
				AppID:              218,
				Name:               "Source SDK Base 2007",
				StateFlags:         4,
				InstallDir:         "Source SDK Base 2007",
				LastUpdated:        1476475691,
				UpdateResult:       0,
				SizeOnDisk:         4046685181,
				BuildID:            "131605",
				BytesToDownload:    2364107440,
				BytesDownloaded:    2364107440,
				AutoUpdateBehavior: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGame, err := tt.i.LoadManifest(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.LoadManifest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotGame, tt.wantGame) {
				t.Errorf("API.LoadManifest() = %v, want %v", gotGame, tt.wantGame)
			}
		})
	}
}
