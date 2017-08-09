package itforarchivists

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/richardlehane/siegfried"
	"github.com/richardlehane/siegfried/pkg/sets"
)

var (
	updateJson      string
	sf              *siegfried.Siegfried
	resultsTemplate *template.Template
	sharedTemplate  *template.Template
	badRequest      = errors.New("bad request")
	thisStore       store
)

func init() {
	// setup global sf
	sf, _ = siegfried.Load("public/latest/pronom-tika-loc.sig")
	// setup global updateJson
	f, _ := os.Open("public/latest/default.sig")
	i, _ := f.Stat()
	current.Size = int(i.Size())
	f.Close()
	s, _ := siegfried.Load("public/latest/default.sig")
	current.Created = s.C.Format(time.RFC3339)
	updateJson = current.Json()

	// templates
	resultsTemplate = template.Must(template.New("resultsT").Parse(templ))
	sharedTemplate = template.Must(template.Must(resultsTemplate.Clone()).Parse(shareTempl))

	// store
	thisStore = make(simpleStore)

	// routes
	http.HandleFunc("/siegfried/identify", hdlErr(handleIdentify))
	http.HandleFunc("/siegfried/update", handleUpdate)
	http.HandleFunc("/siegfried/sets", hdlErr(handleSets))
	http.HandleFunc("/siegfried/results", hdlErr(handleResults))
	http.HandleFunc("/siegfried/results/", hdlErr(handleResults))
	http.HandleFunc("/siegfried/share", hdlErr(handleShare))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Sorry, that doesn't seem to be a valid route :)", 404)
	})
}

func hdlErr(f func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, "Sorry, something went wrong :(: "+err.Error(), 500)
		}
	}
}

func handleIdentify(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		id, err := identify(r)
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(w, id.JSON())
		return nil
	}
	return badRequest
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, updateJson)
}

func handleSets(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			return errors.New("invalid sets form")
		}
		vals, ok := r.Form["set"]
		if !ok {
			return errors.New("invalid sets form")
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(w, "[\""+strings.Join(sets.Sets(vals...), "\", \"")+"\"]")
		return nil
	}
	return badRequest
}

func handleResults(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" && r.URL.Path == "/siegfried/results" {
		return parseResults(w, r)
	}
	if strings.HasPrefix(r.URL.Path, "/siegfried/results/") {
		uuid := strings.TrimPrefix(r.URL.Path, "/siegfried/results/")
		uuid = strings.TrimSuffix(uuid, "/") // remove any trailing slash
		return retrieveResults(w, uuid, thisStore)
	}
	return badRequest
}

func handleShare(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return share(w, r, thisStore)
	}
	return badRequest
}
