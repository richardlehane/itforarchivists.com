package itforarchivists

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/richardlehane/crock32"
	"github.com/richardlehane/runner"
)

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
	u := puuid()
	if err := s.stash(tag+"/"+u, "", title, desc, lg); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, `{"success": "/siegfried/logs/`+tag+`/`+u+`"}`)
	return nil

}

func getLog(rdr io.Reader) (string, string, *runner.Log, error) {
	lg := &runner.Log{}
	dec := json.NewDecoder(rdr)
	if err := dec.Decode(lg); err != nil {
		return "", "", nil, err
	}
	title := lg.Detail
	var desc string
	if len(lg.Reports) > 0 {
		desc = "Log created at " + lg.Reports[0].Start.Format(time.RFC3339)
	}
	return title, desc, lg, nil
}

func retrieveLogs(w http.ResponseWriter, tag, uuid string, s store) error {
	if _, err := crock32.Decode(uuid); err != nil {
		return badRequest
	}
	_, _, _, _, err := s.retrieve(tag + "/" + uuid)
	return err
}

var developJobs = runner.Jobs{
	{
		Detail:   "corpora dir",
		Cmd:      []string{"mkdir", "/root/corpora"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "setup",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
	{
		Detail:   "output dir",
		Cmd:      []string{"mkdir", "/root/out"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "setup",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
	{
		Detail:   "rclone corpora",
		Cmd:      []string{"rclone", "sync", "--transfers=32", "bb:corpora", "/root/corpora"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "setup",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
	{
		Detail: "save profile",
		Background: &runner.Background{
			Delay: 30 * time.Second,
			Cmd:   []string{"sfprof", "-multi", "32", "/root/corpora"},
		},
		Cmd:      []string{"curl", "-vs", "-o", "/root/out/sf.prof", "http://localhost:6060/debug/pprof/profile"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "profile",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
	{
		Detail:   "upload profile",
		Cmd:      []string{"go", "tool", "pprof", "-png", "/root/out/sf.prof"},
		RunTwice: false,
		SendOut:  true,
		Base64:   true,
		LogKey:   "profile",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
	{
		Detail:   "govdocs - master",
		Cmd:      []string{"sf", "-log", "", "-multi", "32", "/root/corpora/govdocs-selected"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "govdocs",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "/root/out/sf_gd.yaml",
	},
	{
		Detail:   "govdocs - develop",
		Cmd:      []string{"sfdev", "-log", "", "-multi", "32", "/root/corpora/govdocs-selected"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "govdocs",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "/root/out/sfdev_gd.yaml",
	},
	{
		Detail:   "govdocs - compare",
		Cmd:      []string{"roy", "compare", "/root/out/sfdev_gd.yaml", "/root/out/sf_gd.yaml"},
		RunTwice: false,
		SendOut:  true,
		Base64:   false,
		LogKey:   "govdocs",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
	{
		Detail:   "ipres - master",
		Cmd:      []string{"sf", "-log", "", "-multi", "32", "/root/corpora/ipres-systems-showcase-files"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "ipres",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "/root/out/sf_ipres.yaml",
	},
	{
		Detail:   "ipres - develop",
		Cmd:      []string{"sfdev", "-log", "", "-multi", "32", "/root/corpora/ipres-systems-showcase-files"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "ipres",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "/root/out/sfdev_ipres.yaml",
	},
	{
		Detail:   "ipres - compare",
		Cmd:      []string{"roy", "compare", "/root/out/sfdev_ipres.yaml", "/root/out/sf_ipres.yaml"},
		RunTwice: false,
		SendOut:  true,
		Base64:   false,
		LogKey:   "ipres",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
	{
		Detail:   "pronom - master",
		Cmd:      []string{"sf", "-log", "", "/root/corpora/pronom-files"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "pronom",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "/root/out/sf_pronom.yaml",
	},
	{
		Detail:   "pronom - develop",
		Cmd:      []string{"sfdev", "-log", "", "/root/corpora/pronom-files"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "pronom",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "/root/out/sfdev_pronom.yaml",
	},
	{
		Detail:   "pronom - compare",
		Cmd:      []string{"roy", "compare", "/root/out/sfdev_pronom.yaml", "/root/out/sf_pronom.yaml"},
		RunTwice: false,
		SendOut:  true,
		Base64:   false,
		LogKey:   "pronom",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
}

var benchJobs = runner.Jobs{
	{
		Detail:   "corpora dir",
		Cmd:      []string{"mkdir", "/root/corpora"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "setup",
		URL:      "https://www.itforarchivists.com/siegfried/logs/bench",
		Save:     "",
	},
	{
		Detail:   "output dir",
		Cmd:      []string{"mkdir", "/root/out"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "setup",
		URL:      "https://www.itforarchivists.com/siegfried/logs/develop",
		Save:     "",
	},
	{
		Detail:   "rclone corpora",
		Cmd:      []string{"rclone", "sync", "--transfers=32", "bb:corpora", "/root/corpora"},
		RunTwice: false,
		SendOut:  false,
		Base64:   false,
		LogKey:   "setup",
		URL:      "https://www.itforarchivists.com/siegfried/logs/bench",
		Save:     "",
	},
	{
		Detail:   "siegfried version",
		Cmd:      []string{"sf", "-version"},
		RunTwice: false,
		SendOut:  true,
		Base64:   false,
		LogKey:   "setup",
		URL:      "https://www.itforarchivists.com/siegfried/logs/bench",
		Save:     "",
	},
}
