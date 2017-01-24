package itforarchivists

import (
	"fmt"
	"net/http"
	"strings"
)

type Identification struct {
	Ids [][][2]string
}

func toJSON(ids [][][2]string) []string {
	ret := make([]string, len(ids))
	for i, id := range ids {
		if id[0] == "namespace" {
			id[0] = "ns"
		}
		ret[i] = fmt.Sprintf("{\"%s\":\"%s\"}", id[0], id[1])
	}
	return ret
}

func (i *Identification) JSON() string {
	return "[" + strings.Join(toJSON(i.Ids), ", ") + "]"
}

func identify(r *http.Request) (*Identification, error) {
	id := &Identification{make([][][2]string, 0, 3)}
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
		id.Ids = append(id.Ids, sf.Label(i))
	}
	return id, nil
}
