package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/richardlehane/crock32"
	"github.com/richardlehane/runner"
)

func retrieveLog(w http.ResponseWriter, path string, s store) error {
	_, _, _, raw, err := s.retrieve(path)
	if err != nil {
		return badRequest
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	buf := bytes.NewBuffer(raw)
	_, err = io.Copy(w, buf)
	return err
}

func parseLogs(w http.ResponseWriter, r *http.Request, tag string, s store) error {
	_, auth, ok := r.BasicAuth()
	if !ok || auth != os.Getenv("RUNNER_AUTH") {
		return fmt.Errorf("bad auth")
	}
	body := r.Body
	defer body.Close()
	title, desc, lg, err := getLog(body)
	if err != nil {
		return err
	}
	path := tag + "/" + crock32.Encode(uint64(lg.Batch.Unix())) + "/" + crock32.PUID()
	if err := s.stash(path, "", title, desc, lg); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, `{"success": "/siegfried/logs/`+path+`"}`)
	return nil

}

func getLog(rdr io.Reader) (string, string, *runner.Log, error) {
	lg := &runner.Log{}
	dec := json.NewDecoder(rdr)
	if err := dec.Decode(lg); err != nil {
		return "", "", nil, err
	}
	title := lg.Label
	var desc string
	if len(lg.Reports) > 0 {
		desc = lg.Batch.Format(time.RFC3339)
	}
	return title, desc, lg, nil
}

func retrieveLogs(w http.ResponseWriter, prefix, dir string, dirs []string, s store) error {
	if _, err := crock32.Decode(dir); err != nil {
		return badRequest
	}
	keys := s.list(prefix, dir)
	if len(keys) == 0 {
		return badRequest
	}
	logs := make([]*runner.Log, len(keys))
	for i, key := range keys {
		_, _, _, raw, err := s.retrieve(key)
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
	return writeLogs(w, prefix, dir, dirs, logs...)
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
		`A data set created for the 2022 iPRES conference Digital Preservation Bake Off Challenge comprising 2,944 files (20.8GB). Includes  includes a number of different content types, ranging from generic and not-so generic PDFs, still images and office documents, to complex objects such as AV, 3D and disk images, to web-based objects such as websites and social media. The ‘exotic ingredients’ section contains data with additional challenges, such as unidentifiable objects, corrupt objects or legacy file formats. Sourced from <a href="https://ipres2022.scot/call-for-contributions-2/data-set/">https://ipres2022.scot/call-for-contributions-2/data-set/</a>.`,
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

var toolDetails = map[string]string{
	"master":  "latest production release",
	"develop": "tip of the <a href='https://github.com/richardlehane/siegfried/tree/develop'>develop branch</a>",
}

var machineDetails = map[string]struct {
	Link        string
	Description string
}{
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
	}, // legacy definitions follow
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

func writeLogs(w http.ResponseWriter, prefix, dir string, dirs []string, logs ...*runner.Log) error {
	var title string
	switch prefix {
	case "bench":
		title = "Siegfried benchmarks"
	case "develop":
		title = "Siegfried development benchmarks"
	}
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
	ts := make([][2]string, len(dirs))
	for i, v := range dirs {
		c, err := crock32.Decode(v)
		if err != nil {
			return err
		}
		ts[i][0] = v
		ts[i][1] = time.Unix(int64(c), 0).String()
	}
	payload := struct {
		Prefix     string
		Title      string
		Time       string
		Machine    Machine
		Profile    string
		Versions   []Tool
		Benchmarks []*Benchmark
		History    [][2]string
	}{
		Prefix:     prefix,
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
