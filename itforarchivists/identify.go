package itforarchivists

import (
	"fmt"
	"net/http"
	"strings"
)

type Identification struct {
	Puids []string
}

func (i *Identification) JSON() string {
	return "[" + strings.Join(i.Puids, ", ") + "]"
}

func identify(r *http.Request) (*Identification, error) {
	id := &Identification{make([]string, 0, 3)}

	// open the submitted file
	f, h, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// identify using the global sf
	c, err := sf.Identify(f, h.Filename, "")
	if err != nil {
		return nil, fmt.Errorf("failed to identify %v, got: %v", h.Filename, err)
	}
	for i := range c {
		id.Puids = append(id.Puids, i.JSON())
	}
	return id, nil
}
