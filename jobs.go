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
	if err := s.stash(tag+"-"+u, "", title, desc, lg); err != nil {
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
	_, _, _, _, err := s.retrieve(tag + "-" + uuid)
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
		Detail:   "siegfried dev version",
		Cmd:      []string{"sfdev", "-version"},
		RunTwice: false,
		SendOut:  true,
		Base64:   false,
		LogKey:   "setup",
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
