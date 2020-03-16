package steampump

import (
	"github.com/m1cr0man/steampump/pkg/steam"
	"github.com/m1cr0man/steampump/pkg/steammesh"
)

type Config struct {
	Steam steam.Config     `json:"steam"`
	Mesh  steammesh.Config `json:"mesh"`
}
