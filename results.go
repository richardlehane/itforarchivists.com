package itforarchivists

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/richardlehane/siegfried/pkg/reader"
)

var fileTitles = []string{"filename", "filesize", "modified", "errors"}

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
	buf := &bytes.Buffer{}
	buf.WriteByte('[')
	for k, v := range c {
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
	//var hh bool
	if len(file.Hash) > 0 {
		//hh = true
		fileTitles = append(fileTitles, head.HashHeader)
	}
	var tool string
	switch {
	case head.Version != [3]int{0, 0, 0}:
		tool = fmt.Sprintf("siegfried v%d.%d.%d", head.Version[0], head.Version[1], head.Version[2])
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
			Titles:     make([]string, len(fileTitles)+len(head.Fields[0])),
			Rows:       make([][]string, 0, 1000),
		}
		copy(results.Datas[i].Titles, fileTitles)
		copy(results.Datas[i].Titles[len(fileTitles):], head.Fields[0])
	}
	for err == nil {
		// DO PROCESSING STUFF HERE
		file, err = res.Next()
	}
	if err == io.EOF {
		err = nil
	}
	return results, err
}
