package main

import (
	"encoding/json"

	"github.com/richardlehane/siegfried/pkg/config"
)

type Update struct {
	Version [3]int `json:"sf"`
	Created string `json:"created"`
	Hash    string `json:"hash"`
	Size    int    `json:"size"`
	Path    string `json:"path"`
}

func (u Update) Json() string {
	byt, _ := json.Marshal(u)
	return string(byt)
}

const domain = "https://www.itforarchivists.com/" // "http://localhost:8081/"

var current = map[string]*Update{
	"pronom": {
		Version: config.Version(),
		Path:    domain + "siegfried/latest",
	},
	"loc": {
		Version: config.Version(),
		Path:    domain + "siegfried/latest/loc",
	},
	"tika": {
		Version: config.Version(),
		Path:    domain + "siegfried/latest/tika",
	},
	"freedesktop": {
		Version: config.Version(),
		Path:    domain + "siegfried/latest/freedesktop",
	},
	"pronom-tika-loc": {
		Version: config.Version(),
		Path:    domain + "siegfried/latest/pronom-tika-loc",
	},
	"deluxe": {
		Version: config.Version(),
		Path:    domain + "siegfried/latest/deluxe",
	},
	"archivematica": {
		Version: config.Version(),
		Path:    domain + "siegfried/latest/archivematica",
	},
}
