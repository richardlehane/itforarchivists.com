package main

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
		idstr := make([]string, len(id))
		if id[0][0] == "namespace" {
			id[0][0] = "ns"
		}
		for j, v := range id {
			idstr[j] = fmt.Sprintf("\"%s\":\"%s\"", v[0], v[1])
		}
		ret[i] = "{" + strings.Join(idstr, ",") + "}"
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
