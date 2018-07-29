package itforarchivists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/richardlehane/crock32"

	"github.com/richardlehane/siegfried/pkg/core"
	"github.com/richardlehane/siegfried/pkg/reader"
)

const noMIME = "no MIME"

var redactFields = []string{"filename"}

func redact(r *Results) *Results {
	for _, d := range r.Datas {
		idxs := make([]int, 0, len(redactFields))
		for i, t := range d.Titles {
			for _, f := range redactFields {
				if t == f {
					idxs = append(idxs, i)
					break
				}
			}
		}
		if len(idxs) > 0 {
			for _, row := range d.Rows {
				for _, idx := range idxs {
					row[idx] = filepath.Ext(row[idx])
				}
			}
		}
	}
	return r
}

func appendUniq(s string, l []string) []string {
	if l == nil {
		return []string{s}
	}
	for _, v := range l {
		if s == v {
			return l
		}
	}
	return append(l, s)
}

func markDupes(r *Results) *Results {
	if len(r.Datas) == 0 {
		return r
	}
	dupesMap := make(map[string][]string)
	for _, row := range r.Datas[0].Rows {
		dupesMap[row[4]] = appendUniq(row[0], dupesMap[row[4]])
	}
	var dupesCount int
	for _, v := range dupesMap {
		if len(v) > 1 {
			dupesCount += len(v)
		}
	}
	if dupesCount == 0 {
		return r
	}
	for _, d := range r.Datas {
		d.Duplicates = dupesCount
		for _, row := range d.Rows {
			if len(dupesMap[row[4]]) > 1 {
				row[len(row)-1] = "true"
			}
		}
	}
	return r
}

func truncate(s string, l int) string {
	if len(s) <= l {
		return s
	}
	return s[:l]
}

func getResults(r *http.Request) (*Results, error) {
	f, _, err := r.FormFile("results")
	if err != nil {
		return nil, err
	}
	dec := json.NewDecoder(f)
	res := &Results{}
	err = dec.Decode(res)
	f.Close()
	if err != nil {
		return nil, fmt.Errorf("bad results: %s", err.Error())
	}
	if !res.validate() {
		return nil, fmt.Errorf("bad results: validation fail")
	}
	return res, nil
}

func share(w http.ResponseWriter, r *http.Request, s store) error {
	name, title, desc, red := r.FormValue("name"), r.FormValue("title"), r.FormValue("description"), r.FormValue("redact")
	name, title, desc = truncate(name, 128), truncate(title, 128), truncate(desc, 256)
	res, err := getResults(r)
	if err != nil {
		return err
	}
	if red != "false" { // redact unless explicitly told not to
		res = redact(res)
	}
	u := crock32.PUID()
	if err := s.stash("results/"+u, name, title, desc, res); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, `{"success": "/siegfried/results/`+u+`"}`)
	return nil
}

func getter(key string, fields [][]string) func(int, core.Identification) (string, int) {
	mm := make(map[int]int)
	for i, v := range fields {
		for j, w := range v {
			switch w {
			case key, strings.ToUpper(key):
				mm[i] = j
			}
		}
	}
	return func(idx int, id core.Identification) (string, int) {
		f, ok := mm[idx]
		if !ok {
			return "", -1
		}
		return id.Values()[f], f
	}
}

var fileTitles = []string{"filename", "filesize", "modified", "errors"}

var hiddenTitles = []string{"hasMultiID", "isDuplicate"}

func (r *Results) validate() bool {
	if len(r.Tool) < 4 {
		return false
	}
	switch r.Tool[:4] {
	case "fido", "droi", "sieg":
	default:
		return false
	}
	if len(r.Datas) == 0 || len(r.Datas) != len(r.Identifiers) {
		return false
	}
	for _, d := range r.Datas {
		t := len(d.Titles)
		if t < len(fileTitles)+1 {
			return false
		}
		for i, ft := range fileTitles {
			if ft != d.Titles[i] {
				return false
			}
		}
		f := len(d.Rows)
		fmts, mimes := d.FmtCounts.total(), d.MIMECounts.total()
		switch {
		case f == 0, f < d.Unknown, f < d.Warn, f < d.Error, f < d.Multiple, f > fmts, f > mimes:
			return false
		case d.Multiple == 0:
			if fmts != mimes || fmts != f {
				return false
			}
		}
		for _, r := range d.Rows {
			if len(r) != t {
				return false
			}
		}
	}
	return true
}

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
	Duplicates int        `json:"duplicates"`
	FmtCounts  Count      `json:"fmtCounts"`
	MIMECounts Count      `json:"mimeCounts"`
	Titles     []string   `json:"titles"`
	Rows       [][]string `json:"rows"`
}

type Count map[string]int

func (c Count) total() int {
	var t int
	for _, v := range c {
		t += v
	}
	return t
}

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

func (c *Count) UnmarshalJSON(byt []byte) error {
	if bytes.Equal(byt, []byte("null")) {
		return nil
	}
	m := make(Count)
	if *c == nil {
		*c = m
	}
	slc := [][2]interface{}{}
	err := json.Unmarshal(byt, &slc)
	if err != nil {
		return err
	}
	for _, row := range slc {
		key, ok := row[0].(string)
		if !ok {
			return fmt.Errorf("bad type in Count, expecting string, got %v", row[0])
		}
		num, ok := row[1].(float64)
		if !ok {
			return fmt.Errorf("bad type in Count, expecting float64, got %T", row[1])
		}
		m[key] = int(num)
	}
	return nil
}

func results(r io.Reader, nm string) (*Results, error) {
	thisFileTitles := make([]string, len(fileTitles))
	copy(thisFileTitles, fileTitles)
	res, err := reader.New(r, nm)
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
		thisFileTitles = append(thisFileTitles, head.HashHeader)
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
			Titles:     make([]string, len(thisFileTitles)+len(head.Fields[i])-1+len(hiddenTitles)),
			Rows:       make([][]string, 0, 1000),
		}
		copy(results.Datas[i].Titles, thisFileTitles)
		copy(results.Datas[i].Titles[len(thisFileTitles):], head.Fields[i][1:]) // skip the first (ns) field
		copy(results.Datas[i].Titles[len(thisFileTitles)+len(head.Fields[i])-1:], hiddenTitles)
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
					last[len(last)-2] = "true"
				}
				multiID = true
			}
			d.FmtCounts[id.String()]++
			mime, _ := mimeGetter(idx, id)
			if mime == "" {
				mime = noMIME
			}
			d.MIMECounts[mime]++
			row := make([]string, len(d.Titles))
			// multiID
			if multiID {
				row[len(row)-2] = "true"
			} else {
				row[len(row)-2] = "false"
			}
			// duplicates
			row[len(row)-1] = "false"
			row[0], row[1], row[2] = file.Path, strconv.FormatInt(file.Size, 10), file.Mod
			if file.Err != nil {
				row[3] = file.Err.Error()
				d.Error++
			}
			if id.Warn() != "" {
				d.Warn++
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
	if hh {
		results = markDupes(results)
	}
	if err == io.EOF {
		err = nil
	}
	return results, err
}

func parseResults(w http.ResponseWriter, r *http.Request) error {
	f, h, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer f.Close()
	res, err := results(f, h.Filename)
	if err != nil {
		return err
	}
	return writeResults(w, res, false, "", "", "", "")
}

// grandfather old results
var oldResults = map[string]bool{
	"13pqzaj": true,
	"396g5jf": true,
	"3hk6wgx": true,
	"959zaj":  true,
	"ea1zaj":  true,
	"wtxzaj":  true,
}

func retrieveResults(w http.ResponseWriter, uuid string, s store) error {
	if _, err := crock32.Decode(uuid); err != nil {
		return badRequest
	}
	key := uuid
	if !oldResults[key] {
		key = "results/" + key
	}
	name, title, desc, raw, err := s.retrieve(key)
	if err != nil {
		return err
	}
	res := &Results{}
	if err := json.Unmarshal(raw, res); err != nil {
		return err
	}
	return writeResults(w, res, true, name, title, desc, uuid)
}

func writeResults(w http.ResponseWriter, res *Results, shared bool, name, title, desc, uuid string) error {
	byt, err := json.Marshal(res)
	if err != nil {
		return err
	}
	templ := resultsTemplate
	if shared {
		templ = sharedTemplate
	}
	return templ.Execute(w,
		struct {
			Name        string
			Title       string
			Desc        string
			UUID        string
			Metadata    [][2]string
			Identifiers []Identifier
			JSON        string
		}{
			Name:  name,
			Title: title,
			Desc:  desc,
			UUID:  uuid,
			Metadata: [][2]string{
				{"Results path", res.ResultsPath},
				{"Tool", res.Tool},
				{"Signature", res.Signature},
				{"Signature created", res.SignatureCreated},
				{"Scan date", res.ScanDate},
			},
			Identifiers: res.Identifiers,
			JSON:        string(byt),
		},
	)
}
