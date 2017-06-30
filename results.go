package itforarchivists

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/richardlehane/siegfried/pkg/core"
	"github.com/richardlehane/siegfried/pkg/reader"
)

func getter(key string, fields [][]string) func(int, core.Identification) string {
	mm := make(map[int]int)
	for i, v := range fields {
		for j, w := range v {
			switch w {
			case key, strings.ToUpper(key):
				mm[i] = j
			}
		}
	}
	return func(idx int, id core.Identification) string {
		f, ok := mm[idx]
		if !ok {
			return ""
		}
		return id.Values()[f]
	}
}

var fileTitles = []string{"filename", "filesize", "modified", "errors"}

var hiddenTitles = []string{"hasWarn", "hasMultiID"}

type Results struct {
	ResultsPath      string       `json:"resultsPath"`
	Tool             string       `json:"tool"`
	Signature        string       `json:"signaturePath"`
	SignatureCreated string       `json:"signatureCreated"`
	ScanDate         string       `json:"scanDate"`
	Identifiers      []Identifier `json:"identifiers"`
	Datas            []*Data      `json:"results"`
}

type Identifier struct {
	Name    string `json:"name"`
	Details string `json:"details"`
}

type Data struct {
	Unknown    int        `json:"unknowns"`
	Warn       int        `json:"warnings"`
	Error      int        `json:"errors"`
	Multiple   int        `json:"multipleIDs"`
	FmtCounts  Count      `json:"fmtCounts"`
	MIMECounts Count      `json:"mimeCounts"`
	Titles     []string   `json:"titles"`
	Rows       [][]string `json:"rows"`
}

type Count map[string]int

func (c Count) MarshalJSON() ([]byte, error) {
	first := true
	buf := &bytes.Buffer{}
	buf.WriteByte('[')
	for k, v := range c {
		if first {
			first = false
		} else {
			fmt.Fprint(buf, ",")
		}
		fmt.Fprintf(buf, "[\"%s\",%d]", k, v)
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

func results(r *http.Request) (*Results, error) {
	f, h, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	res, err := reader.New(f, h.Filename)
	if err != nil {
		return nil, err
	}
	head := res.Head()
	if len(head.Identifiers) < 1 {
		return nil, fmt.Errorf("results file contains no identifiers")
	}
	file, err := res.Next()
	if err != nil {
		return nil, fmt.Errorf("no valid results; got %v", err)
	}
	var hh bool
	if len(file.Hash) > 0 {
		hh = true
		fileTitles = append(fileTitles, head.HashHeader)
	}
	var tool string
	switch {
	case head.Version != [3]int{0, 0, 0}:
		tool = fmt.Sprintf("siegfried %d.%d.%d", head.Version[0], head.Version[1], head.Version[2])
	case head.Identifiers[0][0] == "droid":
		tool = "droid"
	case head.Identifiers[0][0] == "fido":
		tool = "fido"
	}
	results := &Results{
		ResultsPath:      head.ResultsPath,
		Tool:             tool,
		Signature:        head.SignaturePath,
		SignatureCreated: head.Created.Format(time.RFC3339),
		ScanDate:         head.Scanned.Format(time.RFC3339),
	}
	results.Identifiers = make([]Identifier, len(head.Identifiers))
	for i, v := range head.Identifiers {
		results.Identifiers[i].Name = v[0]
		results.Identifiers[i].Details = v[1]
	}
	results.Datas = make([]*Data, len(head.Identifiers))
	for i := range results.Datas {
		results.Datas[i] = &Data{
			FmtCounts:  make(Count),
			MIMECounts: make(Count),
			Titles:     make([]string, len(fileTitles)+len(head.Fields[0])-1+len(hiddenTitles)),
			Rows:       make([][]string, 0, 1000),
		}
		copy(results.Datas[i].Titles, fileTitles)
		copy(results.Datas[i].Titles[len(fileTitles):], head.Fields[0][1:]) // skip the first (ns) field
		copy(results.Datas[i].Titles[len(fileTitles)+len(head.Fields[0])-1:], hiddenTitles)
	}
	mimeGetter := getter("mime", head.Fields)
	for err == nil {
		idx := -1
		var ns string
		var multiID bool
		var d *Data
		for _, id := range file.IDs {
			if id.Values()[0] != ns {
				idx++
				ns = id.Values()[0]
				multiID = false
				d = results.Datas[idx]
			} else {
				if !multiID {
					d.Multiple++
					last := d.Rows[len(d.Rows)-1]
					last[len(last)-1] = "true"
				}
				multiID = true
			}
			d.FmtCounts[id.String()]++
			d.MIMECounts[mimeGetter(idx, id)]++
			row := make([]string, len(d.Titles))
			if multiID {
				row[len(row)-1] = "true"
			} else {
				row[len(row)-1] = "false"
			}
			row[0], row[1], row[2] = file.Path, strconv.FormatInt(file.Size, 10), file.Mod
			if file.Err != nil {
				row[3] = file.Err.Error()
				d.Error++
			}
			if id.Warn() != "" {
				row[len(row)-2] = "true"
				d.Warn++
			} else {
				row[len(row)-2] = "false"
			}
			if !id.Known() {
				d.Unknown++
			}
			slc := 4
			if hh {
				row[4] = string(file.Hash)
				slc = 5
			}
			copy(row[slc:], id.Values()[1:])
			d.Rows = append(d.Rows, row)
		}
		file, err = res.Next()
	}
	if err == io.EOF {
		err = nil
	}
	return results, err
}
