package server

import "time"

type Config struct {
	Port int `json:"port"`
}

type FileInfo struct {
	Name  string    `json="name"`
	Dir   bool      `json="dir"`
	Mtime time.Time `json="mtime"`
}
