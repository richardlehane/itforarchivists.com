package itforarchivists

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/richardlehane/crock32"

	"github.com/richardlehane/siegfried/pkg/core"
	"github.com/richardlehane/siegfried/pkg/reader"
)

var redactFields = []string{"filename"}

const noMIME = "no MIME"

func pseudo(r, t int64) uint64 {
	r = r & 0xFFFF
	t = t << 16 & 0xFFFF0000
	return uint64(t | r)
}

// psuedo uuid
func puuid() string {
	p := pseudo(rand.Int63(), time.Now().UnixNano())
	return crock32.Encode(p)
}

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

func truncate(s string, l int) string {
	if len(s) <= l {
		return s
	}
	return s[:l]
}

func share(w http.ResponseWriter, r *http.Request, s store) error {
	name, title, desc, red := r.FormValue("name"), r.FormValue("title"), r.FormValue("description"), r.FormValue("redact")
	name, title, desc = truncate(name, 128), truncate(title, 128), truncate(desc, 256)
	f, _, err := r.FormFile("results")
	if err != nil {
		return err
	}
	dec := json.NewDecoder(f)
	res := &Results{}
	err = dec.Decode(res)
	f.Close()
	if err != nil {
		return fmt.Errorf("bad results")
	}
	if !res.validate() {
		return fmt.Errorf("bad results")
	}
	switch red {
	case "false", "0":
	default:
		res = redact(res)
	}
	u := puuid()
	if err := s.stash(u, name, title, desc, res); err != nil {
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

var hiddenTitles = []string{"hasMultiID"}

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
			mime, _ := mimeGetter(idx, id)
			if mime == "" {
				mime = noMIME
			}
			d.MIMECounts[mime]++
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
	return writeResults(w, res, false, "", "", "")
}

func retrieveResults(w http.ResponseWriter, uuid string, s store) error {
	if _, err := crock32.Decode(uuid); err != nil {
		return badRequest
	}
	name, title, desc, raw, err := s.retrieve(uuid)
	if err != nil {
		return err
	}
	res := &Results{}
	if err := json.Unmarshal(raw, res); err != nil {
		return err
	}
	return writeResults(w, res, true, name, title, desc)
}

func writeResults(w http.ResponseWriter, res *Results, shared bool, name, title, desc string) error {
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
			Metadata    [][2]string
			Identifiers []Identifier
			JSON        string
		}{
			Name:  name,
			Title: title,
			Desc:  desc,
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

var shareTempl = `{{ define "share" -}} 
	<div class="pure-g pure-u-md-1-4 l-box">
	{{ if len .Title | lt 0 }}<h1>{{ .Title }}</h1>{{ end }}
	{{ if len .Name | lt 0 }}<p><i>{{ .Name }}</i></p>{{ end }}
	{{ if len .Desc | lt 0 }}<p>{{ .Desc }}</p>{{ end }}
	</div>
{{- end }} `

var templ = `
<html>
<head>
<title>{{ if len .Title | eq 0  }}Siegfried results chart{{ else }}{{ .Title }}{{ end }}</title>
<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/1.10.15/css/jquery.dataTables.min.css">
<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/buttons/1.3.1/css/buttons.dataTables.min.css">

<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/pure-min.css" integrity="sha384-nn4HPE8lTHyVtfCBi5yW9d20FjT8BJwUXyWZT9InLYax14RDjBj46LmSztkmNP9w" crossorigin="anonymous">
<link rel="stylesheet" href="https://unpkg.com/purecss@1.0.0/build/grids-responsive-min.css">
<link rel="stylesheet" href="/css/maia.css">

<meta name="viewport" content="width=device-width, initial-scale=1">
    
<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
<script type="text/javascript" src="https://code.jquery.com/jquery-1.12.4.js"></script>
<script type="text/javascript" src="https://cdn.datatables.net/1.10.15/js/jquery.dataTables.min.js"></script>

<script type="text/javascript" src="https://cdn.datatables.net/buttons/1.3.1/js/dataTables.buttons.min.js"></script>
<script type="text/javascript" src="https://cdn.datatables.net/buttons/1.3.1/js/buttons.html5.min.js"></script>
<script type="text/javascript" src="https://cdn.datatables.net/buttons/1.3.1/js/buttons.print.min.js"></script>
<script type="text/javascript">var RESULTS = {{ .JSON }};</script>
<script src="/js/results.js"></script>
<script>
(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
})(window,document,'script','//www.google-analytics.com/analytics.js','ga');
ga('create', 'UA-15845576-1', 'auto');
ga('send', 'pageview');
</script>
<style>
.chart {
	height: 320px;
}
</style>
</head>
<body>
<div class="topcorner">
	<a href="/siegfried">Back to siegfried</a> 
</div>
<div class="pure-g">
{{ block "share" . -}}
	<div class="pure-u-1 pure-u-md-1-4 l-box">
	<h1>Share results</h1>
		<form id="share-form" class="pure-form pure-form-stacked">
		  <fieldset>
			<input id="sharename" type="text" name="name" maxlength="128" size="20" placeholder="Name (or organisation)">
			<input id="sharetitle" type="text" name="title" maxlength="128" size="20" placeholder="Title">
			<textarea id="sharedesc" name="description" maxlength="256" rows="3" cols="20" placeholder="Description"></textarea>
			<label for="redact" class="pure-checkbox">
				<input id="redact" type="checkbox" name="redact" value="true" checked> Redact filenames
			</label>
			<input type="submit" value="Publish" class="pure-button pure-button-primary">
		  </fieldset>
		</form>
	</div>
{{- end }}
<div class="pure-u-1 pure-u-md-1-4 l-box">
<h1>Identifiers</h1>
{{- range $idx, $el := .Identifiers -}}
		<p><a href="#" onclick="load({{ $idx }}); return false;"><strong>{{ $el.Name }}</strong></a><br />{{ $el.Details }}</p>
{{- end -}}
</div>
<div class="pure-u-1 pure-u-md-1-4 l-box">
<h1>Details</h1>
<p>
{{- range $idx, $el := .Metadata }}
	{{- if index $el 1 | len | ne 0 -}}
	  {{- if gt $idx 0 -}}<br />{{ end -}}
	  <strong>{{ index $el 0 }}</strong>: {{ index $el 1 }}
	{{- end -}}
{{ end -}}
</p>
  <p>
    <a href="#" id="errNo"><span></span> errors</a><br />
	<a href="#" id="warnNo"><span></span> warnings</a><br />
	<a href="#" id="unkNo"><span></span> unknowns</a><br />
	<a href="#" id="multiNo"><span></span> multiple IDs</a>
  </p>
  <p>
    <a href="#" id="reset">Reset (show all)</a>
  </p>
</div>
<div class="pure-u-1 pure-u-lg-1-4 l-box" id="charts">
  <div id="fmtchart" class="chart"></div>
  <div id="mimechart" class="chart"></div>
  <div class="pure-button-group centre" role="group">
    <button onclick="reveal('fmtchart'); return false;" class="pure-button pure-button-active">Format IDs</button>
    <button onclick="reveal('mimechart'); return false;" class="pure-button">MIME-types</button>
</div>
</div>
<div class="pure-u-1"><table id="table" class="display pure-table pure-table-bordered" width="100%"></table></div>
</div>
</body>
</html>`
