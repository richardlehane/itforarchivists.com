// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sort"

	"github.com/richardlehane/siegfried/pkg/config"
	"github.com/richardlehane/siegfried/pkg/sets"
)

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

func main() {
	d := &Data{
		Version: version(),
		Keys:    keys(),
	}
	byt, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err = ioutil.WriteFile("data/siegfried.json", byt, 0666); err != nil {
		log.Fatal(err)
	}
	log.Println("Generated new siegfried.json file in the data directory")
}
