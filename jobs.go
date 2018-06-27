package itforarchivists

import (
	"fmt"
	"net/http"
	"os"

	"github.com/richardlehane/crock32"
	"github.com/richardlehane/runner"
)

func parseReports(w http.ResponseWriter, r *http.Request) error {
	auth := r.FormValue("auth")
	if auth != os.Getenv("RUNNER_AUTH") {
		return fmt.Errorf("bad auth token")
	}
	f, _, err := r.FormFile("file")
	if err != nil {
		return err
	}
	return f.Close()
}

func retrieveReports(w http.ResponseWriter, uuid string, s store) error {
	if _, err := crock32.Decode(uuid); err != nil {
		return badRequest
	}
	_, _, _, _, err := s.retrieve(uuid)
	return err
}

var Jobs = runner.Jobs{
	{
		Detail:   "A test job",
		Cmd:      []string{"ls", "-a"},
		RunTwice: false,
		LogKey:   "",
		URL:      "https://www.itforarchivists.com/siegfried/reports",
		Save:     "",
	},
}
