package itforarchivists

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/richardlehane/runner"
	"github.com/richardlehane/siegfried"
	"github.com/richardlehane/siegfried/pkg/sets"
)

var (
	updateJson      map[string]string
	sf              *siegfried.Siegfried
	resultsTemplate *template.Template
	sharedTemplate  *template.Template
	badRequest      = errors.New("bad request")
)

func init() {
	updateJson = make(map[string]string)
	// setup global sf
	sf, _ = siegfried.Load("public/latest/pronom-tika-loc.sig")
	// setup global updateJson
	for k := range current {
		fname := k + ".sig"
		if fname == "pronom.sig" {
			fname = "default.sig"
		}
		f, err := ioutil.ReadFile("public/latest/" + fname)
		if err != nil {
			panic(err)
		}
		current[k].Size = len(f)
		h := sha256.New()
		h.Write(f)
		current[k].Hash = hex.EncodeToString(h.Sum(nil))
		s, err := siegfried.Load("public/latest/" + fname)
		if err != nil {
			panic(err)
		}
		current[k].Created = s.C.Format(time.RFC3339)
		updateJson[k] = current[k].Json()
	}

	// templates
	resultsTemplate = template.Must(template.New("resultsT").Parse(templ))
	sharedTemplate = template.Must(template.Must(resultsTemplate.Clone()).Parse(shareTempl))

	// routes
	http.HandleFunc("/siegfried/identify", hdlErr(handleIdentify))
	http.HandleFunc("/siegfried/update", handleUpdate)
	http.HandleFunc("/siegfried/update/", handleUpdate)
	http.HandleFunc("/siegfried/sets", hdlErr(handleSets))
	http.HandleFunc("/siegfried/results", hdlErr(handleResults))
	http.HandleFunc("/siegfried/results/", hdlErr(handleResults))
	http.HandleFunc("/siegfried/share", hdlErr(handleShare))
	http.HandleFunc("/siegfried/redact", hdlErr(handleRedact))
	http.HandleFunc("/siegfried/jobs/", hdlErr(handleJobs))
	http.HandleFunc("/siegfried/logs/", hdlErr(handleLogs))
	http.HandleFunc("/siegfried/benchmarks", handleBench)
	http.HandleFunc("/siegfried/benchmarks/", handleBench)
	http.HandleFunc("/siegfried/develop", handleDevelop)
	http.HandleFunc("/siegfried/develop/", handleDevelop)
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
	var ret string
	if strings.HasPrefix(r.URL.Path, "/siegfried/update/") {
		sig := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/siegfried/update/"), "/") // remove any trailing slash
		ret = updateJson[sig]
	}
	if ret == "" {
		ret = updateJson["pronom"]
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, ret)
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
		thisStore, err := newCloudStore(r)
		if err != nil {
			return err
		}
		uuid := strings.TrimPrefix(r.URL.Path, "/siegfried/results/")
		uuid = strings.TrimSuffix(uuid, "/") // remove any trailing slash
		return retrieveResults(w, uuid, thisStore)
	}
	return badRequest
}

func handleShare(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		thisStore, err := newCloudStore(r)
		if err != nil {
			return err
		}
		return share(w, r, thisStore)
	}
	return badRequest
}

func handleRedact(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		res, err := getResults(r)
		if err != nil {
			return err
		}
		res = redact(res)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		enc := json.NewEncoder(w)
		return enc.Encode(res)
	}
	return badRequest
}

func handleLogs(w http.ResponseWriter, r *http.Request) error {
	var tag string
	switch {
	case strings.HasPrefix(r.URL.Path, "/siegfried/logs/bench"):
		tag = "bench"
		//uuid = strings.TrimPrefix(r.URL.Path, "/siegfried/logs/bench")
	case strings.HasPrefix(r.URL.Path, "/siegfried/logs/develop"):
		tag = "develop"
		//uuid = strings.TrimPrefix(r.URL.Path, "/siegfried/logs/develop")
	default:
		return badRequest
	}
	thisStore, err := newCloudStore(r)
	if err != nil {
		return err
	}
	if r.Method == "POST" {
		return parseLogs(w, r, tag, thisStore)
	}
	return badRequest
}

func handleJobs(w http.ResponseWriter, r *http.Request) error {
	kind := strings.TrimPrefix(r.URL.Path, "/siegfried/jobs/")
	kind = strings.TrimSuffix(kind, "/") // remove any trailing slash
	var jobs runner.Jobs
	switch kind {
	case "develop":
		jobs = developJobs
	case "bench":
		jobs = benchJobs
	default:
		return badRequest
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	byt, err := json.MarshalIndent(jobs, "", "  ")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, string(byt))
	return err
}

func handleBench(w http.ResponseWriter, r *http.Request) {}

func handleDevelop(w http.ResponseWriter, r *http.Request) {}
