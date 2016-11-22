package itforarchivists

import (
	"fmt"
	"net/http"
	"strings"
)

type Identification struct {
	Ids []string
}

func (i *Identification) JSON() string {
	return "[" + strings.Join(i.Ids, ", ") + "]"
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
	ids, err := sf.Identify(f, h.Filename, "")
	if ids == nil {
		return nil, fmt.Errorf("failed to identify %v, got: %v", h.Filename, err)
	}
	for _, i := range ids {
		id.Ids = append(id.Ids, i.JSON())
	}
	return id, nil
}
