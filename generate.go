// this file is used to generate the pages for migrated results and benchmarks
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/richardlehane/crock32"
	"github.com/richardlehane/runner"
)

const (
	TFMT = "2006-01-02T15:04:05"

	// to refresh these: go to https://datatables.net/download/index, choose Datatables styling, jquery3, Datatables, Buttons-> HTML5 export, CDN -> Minify + Concatentate
	DT_JS            = "https://cdn.datatables.net/v/dt/jq-3.7.0/dt-2.3.8/b-3.2.6/b-html5-3.2.6/datatables.min.js"
	DT_JS_INTEGRITY  = "sha384-c8v8R8HgB2YgkBbz+W+EB5M/1RKD5wqDVq1GWXpVLlYmI2EJ+qhj03f0LWH4nd3E"
	DT_CSS           = "https://cdn.datatables.net/v/dt/jq-3.7.0/dt-2.3.8/b-3.2.6/b-html5-3.2.6/datatables.min.css"
	DT_CSS_INTEGRITY = "sha384-1OScWqbx4XgbHvbfZzPSoQ5TfKTHjWXe0jqxPmrngilLvmJ4Z4F8gkX/b9H/q1z4"
)

var (
	resultsTemplate *template.Template
	indexTemplate   *template.Template
	logsTemplate    *template.Template
)

func main() {
	store, err := newSimpleStore("data")
	if err != nil {
		log.Fatal(err)
	}
	repl := strings.NewReplacer("%DT_CSS%", DT_CSS, "%DT_CSS_INTEGRITY%", DT_CSS_INTEGRITY, "%DT_JS%", DT_JS, "%DT_JS_INTEGRITY%", DT_JS_INTEGRITY)
	rChartCSSTempl = repl.Replace(rChartCSSTempl)
	rChartJSTempl = repl.Replace(rChartJSTempl)
	lCSSTempl = repl.Replace(lCSSTempl)
	lJSTempl = repl.Replace(lJSTempl)
	resultsTemplate = parseTemplate("resultsT", templ, rChartCSSTempl, rChartJSTempl, rDetailsTempl, rContent)
	indexTemplate = parseTemplate("resultsIdx", templ, rTableContentsTempl)
	logsTemplate = parseTemplate("logsT", templ, lCSSTempl, lJSTempl, lContent)
	// write results
	dirs := store.results()
	out, err := os.Create(filepath.Join("assets", "attic", "results", "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	store.writeResultsIndex(out, dirs)
	out.Close()
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join("assets", "attic", "results", dir), 0777)
		if err != nil {
			log.Fatal(err)
		}
		out, err = os.Create(filepath.Join("assets", "attic", "results", dir, "index.html"))
		if err != nil {
			log.Fatal(err)
		}
		store.writeResult(out, dir)
		out.Close()
	}
	// write benchmarks
	dirs = store.benchmarks()
	out, err = os.Create(filepath.Join("assets", "attic", "benchmarks", "index.html"))
	if err != nil {
		log.Fatal(err)
	}
	store.writeBenchIndex(out, dirs)
	out.Close()
	for _, dir := range dirs {
		out, err = os.Create(filepath.Join("assets", "attic", "benchmarks", dir, "index.html"))
		if err != nil {
			log.Fatal(err)
		}
		store.writeBench(out, dir, dirs)
		out.Close()
	}
}

type Item struct {
	Path     string `json:"name"`
	Creation string `json:"creation_time"`
	Custom   struct {
		Title       string `json:"x-goog-meta-title"`
		Name        string `json:"x-goog-meta-name"`
		Description string `json:"x-goog-meta-description"`
	} `json:"custom_fields"`
}

type StoreItem struct {
	Name     string
	Creation time.Time
	Title    string
	Desc     string
	Dir      string
}

type simpleStore struct {
	m    map[string]StoreItem
	base string
}

func (ss *simpleStore) load(dir string) error {
	var items []Item
	byt, err := os.ReadFile(filepath.Join(ss.base, dir, "index.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(byt, &items)
	if err != nil {
		return err
	}
	for _, i := range items {
		if _, ok := ss.m[i.Path]; ok {
			return errors.New("duplicate key!")
		}
		t, err := time.Parse(TFMT, i.Creation[:len(i.Creation)-5])
		if err != nil {
			return err
		}
		ss.m[i.Path] = StoreItem{
			Name:     i.Custom.Name,
			Creation: t,
			Title:    i.Custom.Title,
			Desc:     i.Custom.Description,
			Dir:      dir,
		}
	}
	return nil
}

func newSimpleStore(path string) (*simpleStore, error) {
	ss := &simpleStore{
		base: path,
		m:    make(map[string]StoreItem),
	}
	if err := ss.load("benchmarks"); err != nil {
		return nil, err
	}
	if err := ss.load("results"); err != nil {
		return nil, err
	}
	return ss, nil
}

func (ss *simpleStore) retrieve(key string) (string, string, string, []byte, error) {
	v, ok := ss.m[key]
	if !ok {
		return "", "", "", nil, fmt.Errorf("can't find key: %s", key)
	}
	key = strings.ReplaceAll(key, "/", string(filepath.Separator))
	path := filepath.Join(ss.base, v.Dir, key)
	byt, err := os.ReadFile(path)
	return v.Name, v.Title, v.Desc, byt, err
}

func (ss *simpleStore) list(prefix string) []string {
	ret := make([]string, 0, 10)
	for k := range ss.m {
		if strings.HasPrefix(k, prefix) {
			ret = append(ret, k)
		}
	}
	return ret
}

func (ss *simpleStore) benchmarks() []string {
	uniques := make(map[string]struct{})
	for k := range ss.m {
		if strings.ContainsRune(k, '/') {
			uniques[k[:strings.IndexRune(k, '/')]] = struct{}{}
		}
	}
	ret := make([]string, 0, len(uniques))
	for k := range uniques {
		ret = append(ret, k)
	}
	sort.Sort(sort.Reverse(crock32.Sortable(ret)))
	return ret
}

type sortable struct {
	ss   *simpleStore
	strs []string
}

func (s *sortable) Len() int {
	return len(s.strs)
}

func (s *sortable) Less(i, j int) bool {
	return s.ss.m[s.strs[i]].Creation.Before(s.ss.m[s.strs[j]].Creation)
}

func (s *sortable) Swap(i, j int) {
	swap := s.strs[i]
	s.strs[i] = s.strs[j]
	s.strs[j] = swap
}

func (ss *simpleStore) results() []string {
	ret := make([]string, 0, 20)
	for k := range ss.m {
		if !strings.ContainsRune(k, '/') {
			ret = append(ret, k)
		}
	}
	srt := &sortable{ss: ss, strs: ret}
	sort.Sort(sort.Reverse(srt))
	return srt.strs
}

// RESULTS

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
	slc := [][2]any{}
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

func (s *simpleStore) writeResultsIndex(w io.Writer, entries []string) error {
	items := make([]StoreItem, len(entries))
	for i, e := range entries {
		items[i] = s.m[e]
	}
	return indexTemplate.Execute(w, struct {
		Title   string
		Desc    string
		Entries []string
		Items   []StoreItem
	}{
		Title:   "Old Results Reports",
		Desc:    `The "Chart your results" feature allowed users to upload their siegfried results for analysis. The feature is discontinued. Legacy results are still available:`,
		Entries: entries,
		Items:   items,
	})
}

func (s *simpleStore) writeResult(w io.Writer, uuid string) error {
	name, title, desc, raw, err := s.retrieve(uuid)
	if err != nil {
		return err
	}
	res := &Results{}
	if err := json.Unmarshal(raw, res); err != nil {
		return err
	}
	return writeResults(w, res, name, title, desc, uuid)
}

func writeResults(w io.Writer, res *Results, name, title, desc, uuid string) error {
	byt, err := json.Marshal(res)
	if err != nil {
		return err
	}
	return resultsTemplate.Execute(w,
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

// Benchmarks

func (s *simpleStore) writeBenchIndex(w io.Writer, entries []string) error {
	items := make([]StoreItem, len(entries))
	for i, e := range entries {
		var it StoreItem
		for k, v := range s.m {
			if strings.HasPrefix(k, e) {
				v.Title = ""
				it = v
				break
			}
		}
		items[i] = it
	}
	return indexTemplate.Execute(w, struct {
		Title   string
		Desc    string
		Entries []string
		Items   []StoreItem
	}{
		Title:   "Old Benchmarks",
		Desc:    "Automated benchmarks are no longer run for siegfried. The legacy production benchmarks are:",
		Entries: entries,
		Items:   items,
	})
}

type Benchmark struct {
	Title       string
	Description string
	Src         string
	Tools       []Tool
	CompareDesc string
	CompareHdrs []string
	Compare     [][]string
}

type Machine struct {
	Label       string
	Link        string
	Description string
}

type Tool struct {
	Label       string
	Version     string
	Description string
	Duration    string
}

var benchDetails = map[string]struct {
	Title       string
	Description string
}{
	"govdocs": {
		"Govdocs (Selected)",
		`A selection from the Govdocs1 corpus comprising 26,124 files (31.4GB). Represents typical office formats, including approx. 15,000 PDFs. Originally sourced from <a href="http://openpreservation.org/blog/2012/07/26/1-million-21000-reducing-govdocs-significantly/">http://openpreservation.org/blog/2012/07/26/1-million-21000-reducing-govdocs-significantly/</a>.`,
	},
	"ipres2022": {
		"iPRES 2022 Pantry",
		`A corpus created for the 2022 iPRES conference Digital Preservation Bake Off Challenge comprising 2,944 files (20.8GB). Includes a number of different content types, ranging from generic and not-so generic PDFs, still images and office documents, to complex objects such as AV, 3D and disk images, to web-based objects such as websites and social media. The ‘exotic ingredients’ section contains data with additional challenges, such as unidentifiable objects, corrupt objects or legacy file formats. Sourced from <a href="https://ipres2022.scot/call-for-contributions-2/data-set/">https://ipres2022.scot/call-for-contributions-2/data-set/</a>.`,
	},
	"ipres": {
		"iPRES Systems Showcase",
		`A corpus created for the 2014 iPRES conference comprising 2,206 files (5GB). Represents a range of formats, including AV and some uncommon types. Sourced from <a href="http://www.webarchive.org.uk/datasets/ipres.ds.1/">http://www.webarchive.org.uk/datasets/ipres.ds.1/</a>.`,
	},
	"pronom": {
		"PRONOM files",
		"A corpus created by Greg Lepore and comprising 1,205 files (2.1GB). Includes a single sample of as many of the PRONOM IDs (PUIDs) that Greg could find.",
	},
	"deluxe": {
		"The Deluxe",
		"This benchmark checks multi-ID identification using the deluxe.sig signature file which contains four identifiers: PRONOM, LOC FDDs, freedesktop.org and tika-mimetypes. This benchmark is run against the PRONOM files corpus.",
	},
	"zip": {
		"Unzipping",
		"This benchmark checks the `sf -z` command (scans within zip files and other container formats) when run against the iPres corpus.",
	},
}

var machineDetails = map[string]struct {
	Link        string
	Description string
}{
	// cherry servers
	"e3_1240v3": {
		Link:        "https://www.cherryservers.com/pricing/dedicated-servers/e3_1240v3?b=37&r=1",
		Description: "4 cores @ 3.4GHz, 16GB ECC DDR3 RAM, 2x SSD 250GB",
	},
	"e3_1240v5": {
		Link:        "https://www.cherryservers.com/pricing/dedicated-servers/e5_1620v4?b=37&r=1",
		Description: "4 cores @ 3.5GHz, 32GB ECC DDR4 RAM, 2x SSD 250GB",
	},
	"e3_1240lv5": {
		Link:        "https://www.cherryservers.com/pricing/dedicated-servers/e3_1240v5?b=37&r=1",
		Description: "4 cores @ 3.5GHz, 32GB ECC DDR4 RAM, 2x SSD 250GB",
	},
	"e5_1620v4": {
		Link:        "https://www.cherryservers.com/pricing/dedicated-servers/e3_1240lv5?b=37&r=1",
		Description: "4 cores @ 2.1GHz, 32GB ECC DDR4 RAM, 2x SSD 250G",
	},
	// equinix metal servers
	"c3.small.x86": {
		Link:        "https://metal.equinix.com/product/servers/c3-small/",
		Description: "8 cores @ 3.40 GHz, 32GB RAM, 960 GB SSD",
	},
	"c3.medium.x86": {
		Link:        "https://metal.equinix.com/product/servers/c3-medium/",
		Description: "24 cores @ 2.8 GHz, 64GB DDR4 RAM, 960 GB SSD",
	},
	"m3.small.x86": {
		Link:        "https://metal.equinix.com/product/servers/m3-small/",
		Description: "8 cores @ 2.8 GHz, 64GB RAM, 960 GB SSD",
	},
	"m3.large.x86": {
		Link:        "https://metal.equinix.com/product/servers/m3-large/",
		Description: "32 cores @ 2.5 GHz, 256GB DDR4 RAM, 2 x 3.8 TB NVMe",
	},
	"s3.xlarge.x86": {
		Link:        "https://metal.equinix.com/product/servers/s3-xlarge/",
		Description: "24 cores @ 2.2 GHz, 192GB DDR4 RAM, 1.9 TB SSD",
	},
	// legacy servers
	"c2.medium.x86": {
		Link:        "https://www.packet.net/bare-metal/servers/c2-medium-epyc/",
		Description: "24 Physical Cores @ 2.2 GHz; 64 GB ECC RAM; 960 GB SSD",
	},
	"m2.xlarge.x86": {
		Link:        "https://www.packet.net/bare-metal/servers/m2-xlarge/",
		Description: "28 Physical Cores @ 2.2 GHz; 384 GB DDR4 ECC RAM; 3.8 TB NVMe Flash",
	},
	"Standard": {
		Link:        "https://www.packet.net/bare-metal/services/storage/",
		Description: "Elastic block storage",
	},
	"Performance": {
		Link:        "https://www.packet.net/bare-metal/services/storage/",
		Description: "Elastic block storage",
	},
	"c1.large.arm": {
		Link:        "https://www.packet.net/bare-metal/servers/c1-large-arm/",
		Description: "96 Physical Cores @ 2.0 GHz; 128 GB DDR4 ECC RAM; 250 GB SSD",
	},
	"c1.small.x86": {
		Link:        "https://www.packet.net/bare-metal/servers/c1-small/",
		Description: "4 Physical Cores @ 3.5 GHz; 32 GB DDR3 ECC RAM; 120 GB SSD",
	},
	"c1.xlarge.x86": {
		Link:        "https://www.packet.net/bare-metal/servers/m2-xlarge/",
		Description: "16 Physical Cores @ 2.6 GHz; 128 GB DDR4 ECC RAM; 1.6 TB NVMe Flash",
	},
	"m1.xlarge.x86": {
		Link:        "https://www.packet.net/bare-metal/servers/m1-xlarge/",
		Description: "24 Physical Cores @ 2.2 GHz; 256 GB DDR4 ECC RAM; 2.8 TB SSD",
	},
	"s1.large.x86": {
		Link:        "https://www.packet.net/bare-metal/servers/s1-large/",
		Description: "16 Physical Cores @ 2.1 GHz; 128 GB DDR4 ECC RAM; 960 GB SSD; 24 TB HDD",
	},
	"t1.small.x86": {
		Link:        "https://www.packet.net/bare-metal/servers/t1-small/",
		Description: "4 Physical Cores @ 2.4 GHz; 8 GB DDR3 RAM; 80 GB SSD",
	},
	"x1.small.x86": {
		Link:        "https://www.packet.net/bare-metal/servers/x1-small/",
		Description: "4 Physical Cores @ 2.0 GHz; 32 GB DDR4-2400 ECC RAM; 240 GB SSD",
	},
}

func benchDetail(detail string) (bool, []string) {
	if strings.HasPrefix(detail, "compare") {
		bits := strings.SplitN(detail, " - ", 2)
		if len(bits) < 2 {
			return true, nil
		}
		return true, strings.Split(bits[1], ", ")
	}
	return false, nil
}

func compressCompare(recs [][]string) [][]string {
	crpl := strings.NewReplacer(";", "; ")
	ret := make([][]string, 0, len(recs))
	for _, rec := range recs {
		var omit bool
		if len(rec) > 3 {
			xfmt111 := true
			omit = true
			for idx, val := range rec[1:] {
				if val == "x-fmt/111" && xfmt111 {
				} else if val == "UNKNOWN" {
					xfmt111 = false
				} else {
					omit = false
				}
				rec[idx+1] = crpl.Replace(val)
			}
		}
		if !omit {
			ret = append(ret, rec)
		}
	}
	return ret
}

func toBench(l *runner.Log, tools []Tool) *Benchmark {
	b := &Benchmark{
		Title:       benchDetails[l.Label].Title,
		Description: benchDetails[l.Label].Description,
		Src:         l.Path,
	}
	for _, rep := range l.Reports {
		cmp, hdrs := benchDetail(rep.Detail)
		if cmp {
			if strings.HasPrefix(rep.Output, "COMPLETE MATCH") {
				b.Compare = [][]string{}
			} else if strings.Contains(rep.Output, ",") {
				recs, _ := csv.NewReader(strings.NewReader(rep.Output)).ReadAll()
				b.Compare = recs
				b.CompareHdrs = hdrs
			}
		} else {
			for _, tool := range tools {
				if tool.Label == rep.Detail {
					tool.Duration = rep.Duration.String()
					b.Tools = append(b.Tools, tool)
				}
			}
		}
	}
	if b.Compare == nil {
		b.CompareDesc = "One or more of the tools failed, so a comparison is not possible."
		return b
	}
	b.CompareDesc = fmt.Sprintf("The tools differed in output for %d files in the corpus.", len(b.Compare))
	red := compressCompare(b.Compare)
	dif := len(b.Compare) - len(red)
	if dif > 0 {
		b.CompareDesc += fmt.Sprintf(" %d of those differences are because of siegfried's use of a text identification algorithm and the following chart excludes those files.", dif)
	}
	b.Compare = red
	return b
}

func toInfo(l *runner.Log) (Machine, []Tool, []Tool) {
	machine := Machine{
		Label: l.Machine,
	}
	if det, ok := machineDetails[l.Machine]; ok {
		machine.Link = det.Link
		machine.Description = det.Description
	}
	var tools []Tool
	var versions []Tool
	for _, rep := range l.Reports {
		if strings.HasPrefix(rep.Detail, "info") {
			bits := strings.SplitN(rep.Detail, " - ", 3)
			if len(bits) == 3 {
				tools = append(tools, Tool{
					Label:       bits[1],
					Description: bits[2],
				})
			}
			continue
		}
		if strings.HasPrefix(rep.Detail, "version") {
			bits := strings.SplitN(rep.Detail, " - ", 2)
			if len(bits) == 2 {
				var present bool
				vsn := strings.TrimSpace(rep.Output)
				for i, v := range versions {
					if v.Label == bits[1] {
						versions[i].Version = v.Version + " " + vsn
						present = true
						break
					}
				}
				if !present {
					versions = append(versions, Tool{
						Label:   bits[1],
						Version: vsn,
					})
				}
			}
			continue
		}
	}
	return machine, tools, versions
}

func (ss *simpleStore) writeBench(w io.Writer, dir string, history []string) error {
	keys := ss.list(dir)
	logs := make([]*runner.Log, len(keys))
	for i, key := range keys {
		_, _, _, raw, err := ss.retrieve(key)
		if err != nil {
			return err
		}
		log := &runner.Log{}
		if err := json.Unmarshal(raw, log); err != nil {
			return err
		}
		log.Path = key
		logs[i] = log
	}
	return writeLogs(w, dir, history, logs...)
}

func writeLogs(w io.Writer, dir string, history []string, logs ...*runner.Log) error {
	const title = "Siegfried benchmarks"
	var profile string
	var machine Machine
	var tools []Tool
	var versions []Tool
	sort.Slice(logs, func(i, j int) bool {
		if len(logs[i].Reports) == 0 {
			return true
		}
		if len(logs[j].Reports) == 0 {
			return false
		}
		return logs[i].Reports[0].Start.Before(logs[j].Reports[0].Start)
	})
	var benchmarks []*Benchmark
	for _, l := range logs {
		switch l.Label {
		default:
			benchmarks = append(benchmarks, toBench(l, tools))
		case "setup":
			machine, tools, versions = toInfo(l)
		case "profile":
			if len(l.Reports) > 1 {
				profile = l.Reports[1].Output
			}
		}
	}
	c, err := crock32.Decode(dir)
	if err != nil {
		return err
	}
	t := time.Unix(int64(c), 0)
	ts := make([][2]string, len(history))
	for i, v := range history {
		c, err := crock32.Decode(v)
		if err != nil {
			return err
		}
		ts[i][0] = v
		ts[i][1] = time.Unix(int64(c), 0).String()
	}
	payload := struct {
		Title      string
		Time       string
		Machine    Machine
		Profile    string
		Versions   []Tool
		Benchmarks []*Benchmark
		History    [][2]string
	}{
		Title:      title,
		Time:       t.Format(time.RFC1123),
		Machine:    machine,
		Profile:    profile,
		Versions:   versions,
		Benchmarks: benchmarks,
		History:    ts,
	}
	return logsTemplate.Execute(w, payload)
}

// TEMPlATES

func parseTemplate(name string, templs ...string) *template.Template {
	t := template.New(name)
	for _, templ := range templs {
		t = template.Must(t.Parse(templ))
	}
	return t
}

// Main template
var templ = `
<!DOCTYPE html>
<html>
	<head id="head">
		<meta charset="utf-8">
		<meta name="viewport" content="width=device-width">
		<meta name="description" content="IT for archivists, siegfried format identification tool">
		<title>IT for archivists</title>
		{{ block "incCSS" . }}{{ end }}
		<link type="text/css" rel="stylesheet" href="/style.css">
		<link rel="icon" type="image/png" href="/img/richard.png">
		{{ block "incJS" . }}{{ end }}
	</head>
	<body class="full-width">
	<header>
		<div class="nav">
			<ul>
				<li><a class="nav-tab" href="/">Home</a></li>
				<li><a class="nav-tab" href="/siegfried">Siegfried</a></li>
				<li><a class="nav-tab" href="/attic">Attic</a></li>
			</ul>
		</div>
	</header>
	<div>
		{{ block "content" . }}{{ end }}
	</div>
	</body>
</html>
`

var rTableContentsTempl = `{{ define "content" -}}
<h1>{{ .Title }}</h1>
<p>{{ .Desc }}</p>
{{ $items := .Items }}
{{- range $idx, $el := .Entries -}}
{{- with $item := index $items $idx}}
<p><a href="{{ $el }}">{{ $el }}</a>{{ if $item.Name }} {{ $item.Name }}{{ end }}{{ if $item.Title }} {{ $item.Title }}{{ end }} ({{ $item.Creation.Format "2 Jan 2006" }})</p>
{{- end -}}
{{- end -}}
{{- end }} `

var rChartCSSTempl = `{{ define "incCSS" -}} 
<link rel="stylesheet" href="%DT_CSS%" integrity="%DT_CSS_INTEGRITY%" crossorigin="anonymous">
<style>
.chart {
	height: 320px;
}
</style>
{{- end }} `

var rChartJSTempl = `{{ define "incJS" -}} 
<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>
<script type="text/javascript" src="%DT_JS%" integrity="%DT_JS_INTEGRITY%" crossorigin="anonymous"></script>
<script type="text/javascript">var RESULTS = {{ .JSON }};</script>
<script src="/attic/results/results.js"></script>
{{- end }} `

var rDetailsTempl = `{{ define "details" -}} 
<div class="chart-box">
	<h1>{{ if len .Title | lt 0 }}{{ .Title }}{{ else }}Untitled{{ end }}</h1>
	<p><i>{{ if len .Name | lt 0 }}{{ .Name }}{{ end }}</i></p>
	<p>{{ if len .Desc | lt 0 }}{{ .Desc }}{{ end }}</p>
</div>
{{- end }} `

var rContent = `{{ define "content" -}}
<div class="chart-container"> 
{{ block "details" . }}{{ end }}
<div class="chart-box">
<h1>Identifiers</h1>
{{- range $idx, $el := .Identifiers -}}
		<p><a href="#" onclick="load({{ $idx }}); return false;"><strong>{{ $el.Name }}</strong></a><br>{{ $el.Details }}</p>
{{- end -}}
</div>
<div class="chart-box">
<h1>Details</h1>
<p>
{{- range $idx, $el := .Metadata }}
	{{- if index $el 1 | len | ne 0 -}}
	  {{- if gt $idx 0 -}}<br>{{ end -}}
	  <strong>{{ index $el 0 }}</strong>: {{ index $el 1 }}
	{{- end -}}
{{ end -}}
</p>
  <p>
    <a href="#" id="errNo"><span></span> errors</a><br>
	<a href="#" id="warnNo"><span></span> warnings</a><br>
	<a href="#" id="unkNo"><span></span> unknowns</a><br>
	<a href="#" id="multiNo"><span></span> multiple IDs</a><br>
	<a href="#" id="dupesNo"><span></span> duplicates</a>
  </p>
  <p>
    <a href="#" id="reset">Reset (show all)</a>
  </p>
</div>
<div class="chart-box" id="charts">
  <div id="fmtchart" class="chart"></div>
  <div id="mimechart" class="chart"></div>
  <div class="centre" role="group">
    <button onclick="reveal('fmtchart'); return false;">Format IDs</button>
    <button onclick="reveal('mimechart'); return false;">MIME-types</button>
  </div>
</div>
</div>
<br>
<table id="table" class="display"></table>
{{- end }} `

// Log templates

var lCSSTempl = `{{ define "incCSS" -}} 
<link rel="stylesheet" href="%DT_CSS%" integrity="%DT_CSS_INTEGRITY%" crossorigin="anonymous">
{{- end -}}
`
var lJSTempl = `{{ define "incJS" -}} 
<script type="text/javascript" src="%DT_JS%" integrity="%DT_JS_INTEGRITY%" crossorigin="anonymous"></script>
{{- end -}}
`

var lContent = `{{ define "content" -}} 
<h1>{{ .Title }}</h1>
<h2>{{ .Time }}</h2>
<h3>Environment</h3>
<p>These benchmarks were <a href="https://github.com/richardlehane/runner">run</a> on a <a href="{{.Machine.Link}}">{{ .Machine.Label}}</a> machine that was <a href="https://github.com/richardlehane/provisioner">automatically provisioned</a>.</p>
<p>Specs for the <a href="{{.Machine.Link}}">{{ .Machine.Label}}</a>: {{.Machine.Description}}.</p>
{{ if len .Versions | lt 0 -}}
<table>
	<caption>List of tools benchmarked</caption>
	<thead>
		<tr>
			<th>Tool</th>
			<th>Version</th>
		</tr>
	</thead>
	<tbody>
		{{- range .Versions -}}
		<tr>
			<td>{{ .Label }}</td>
			<td>{{ .Version }}</td>
		</tr>
		{{- end -}}
	</tbody>
</table>
{{- end -}}
{{- range $idx, $el := .Benchmarks -}}
<div>
<div class="page-header">
<h2>{{ .Title }}</h2>
</div>
<p>{{ .Description }}</p>
<table>
	<caption>Results</caption>
	<thead>
		<tr>
			<th>Tool</th>
			<th>Description</th>
			<th>Duration</th>
		</tr>
	</thead>
	<tbody>
		{{- range .Tools -}}
		<tr>
			<td>{{ .Label }}</td>
			<td>{{ .Description }}</td>
			<td>{{ .Duration }}</td>
		</tr>
		{{- end -}}
	</tbody>
</table>
<p>{{ .CompareDesc }}</p>
{{ if len .Compare | lt 0 -}}
<table id="cmp{{ $idx }}">
	<caption>Differences between tools</caption>
	<thead>
		<tr>
		<td>file</td>
		{{- range .CompareHdrs -}}
		<td>{{ . }}</td>
		{{- end -}}
		</tr>
	</thead>
	<tbody>
		{{- range $row := .Compare -}}
		<tr>
		{{- range $row -}}
			<td>{{ . }}</td>
		{{- end -}}
		</tr>
		{{- end -}}
	</tbody>
</table>
<script>
$(document).ready(function() {
    $('#cmp{{ $idx }}').DataTable();
} );
</script>
{{ end }}
<p><a href="/attic/benchmarks/{{ .Src }}">Raw output</a></p>
</div>			
{{- end -}}
<div>
<h2>History</h2>
{{- range .History -}}
<p><a href="/attic/benchmarks/{{ index . 0 }}">{{ index . 1 }}</a></p>
{{- end -}}
</div>
{{- end -}}
`
