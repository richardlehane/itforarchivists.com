package itforarchivists

import (
	"encoding/json"
)

type Update struct {
	Version [3]int `json:"sf"`
	Created string `json:"created"`
	Size    int    `json:"size"`
	Path    string `json:"path"`
}

func (u Update) Json() string {
	byt, _ := json.Marshal(u)
	return string(byt)
}
