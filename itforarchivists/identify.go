package itforarchivists

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/richardlehane/siegfried"
)

type Identification struct {
	Puids []string
}

func (i *Identification) JSON() string {
	return "[" + strings.Join(i.Puids, ", ") + "]"
}

func identify(r *http.Request) (*Identification, error) {
	s, err := siegfried.Load("public/latest/pronom-tika-loc.sig")
	id := &Identification{}
	f, h, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	c, err := s.Identify(f, h.Filename, "")
	if err != nil {
		return nil, fmt.Errorf("failed to identify %v, got: %v", h.Filename, err)
	}
	for i := range c {
		id.Puids = append(id.Puids, i.JSON())
	}
	return id, nil
}
