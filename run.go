//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/richardlehane/siegfried/pkg/config"
	"github.com/richardlehane/siegfried/pkg/sets"
)

func main() {
	_ = config.Home() // stopgap due to bug in 1.11.1
	d := &Data{
		Version: version(),
		Keys:    keys(),
	}
	byt, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err = os.WriteFile(filepath.Join("data", "siegfried.json"), byt, 0666); err != nil {
		log.Fatal(err)
	}
	log.Println("Generated new siegfried.json file in the data directory")
	// 1.10.0 update
	wf1 := func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) != ".sig" {
			return nil
		}
		u, err := makeUpdate(path)
		if err != nil {
			return err
		}
		fn := strings.TrimSuffix(filepath.Base(path), ".sig")
		u.Path = domain + "siegfried/latest/1_10/" + fn
		byt, err := json.Marshal(u)
		if err != nil {
			return err
		}
		if fn == "default" {
			fn = "update"
		}
		if err = os.WriteFile(filepath.Join("static", "update", fn+".json"), byt, 0666); err != nil {
			return err
		}
		return nil
	}
	if err := filepath.Walk(filepath.Join("static", "latest", "1_10"), wf1); err != nil {
		log.Fatal(err)
	}
	log.Println("Generated new update json files for 1.10.x")
	// >= 1.11
	updateVersions := make(map[string]Updates)
	wf2 := func(path string, info fs.FileInfo, err error) error {
		bits := strings.Split(path, string(filepath.Separator))
		if len(bits) != 4 || filepath.Ext(bits[3]) != ".sig" {
			return nil
		}
		if bits[2] == "1_10" {
			return nil
		}
		fn := strings.TrimSuffix(bits[3], ".sig")
		u, err := makeUpdate(path)
		if err != nil {
			return err
		}
		u.Path = domain + "siegfried/latest/" + bits[2] + "/" + fn
		if fn == "default" {
			fn = "update"
		}
		updateVersions[fn] = append(updateVersions[fn], u)
		return nil
	}
	if err := filepath.Walk(filepath.Join("static", "latest"), wf2); err != nil {
		log.Fatal(err)
	}
	for k, v := range updateVersions {
		byt, err := json.Marshal(v)
		if err != nil {
			log.Fatal(err)
		}
		if err = os.WriteFile(filepath.Join("static", "update", "v2", k+".json"), byt, 0666); err != nil {
			log.Fatal(err)
		}
	}
	log.Println("Generated new update json files for >= 1.11.x")
}

// generate data/siegfried.json

const templ = "https://github.com/richardlehane/siegfried/releases/download/v"

type Data struct {
	Version string
	Keys    []string
}

func version() string {
	v := config.Version()
	return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2])
}

func keys() []string {
	k := sets.Keys()
	keys := make([]string, 0, len(k))
	for _, v := range k {
		if v[0] > 57 || v[0] < 46 {
			keys = append(keys, v)
		}
	}
	sort.Strings(keys)
	return keys
}

// generate update.json files

type Updates []Update

type Update struct {
	Version [3]int `json:"sf"`
	Created string `json:"created"`
	Hash    string `json:"hash"`
	Size    int    `json:"size"`
	Path    string `json:"path"`
}

func (u Update) Json() string {
	byt, _ := json.Marshal(u)
	return string(byt)
}

const domain = "https://www.itforarchivists.com/" // "http://localhost:8081/"

func makeUpdate(path string) (Update, error) {
	var u Update
	fbuf, err := os.ReadFile(path)
	if err != nil {
		return u, err
	}
	u.Size = len(fbuf)
	h := sha256.New()
	h.Write(fbuf)
	u.Hash = hex.EncodeToString(h.Sum(nil))
	if len(fbuf) < len(config.Magic())+2+15 {
		return u, errors.New("invalid signature file " + path)
	}
	u.Version[0] = int(fbuf[len(config.Magic())])
	u.Version[1] = int(fbuf[len(config.Magic())+1])
	rc := flate.NewReader(bytes.NewBuffer(fbuf[len(config.Magic())+2:]))
	nbuf := make([]byte, 15)
	if n, err := rc.Read(nbuf); n < 15 {
		return u, err
	}
	rc.Close()
	t := &time.Time{}
	if err := t.UnmarshalBinary(nbuf); err != nil {
		return u, err
	}
	u.Created = t.Format(time.RFC3339)
	return u, nil
}
