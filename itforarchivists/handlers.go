package itforarchivists

import (
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/richardlehane/pronom/sets"
)

func wrapError(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func serve404(w http.ResponseWriter) error {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return T("404").Execute(w, nil)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, current.Json())
}

func handleIdentify(w http.ResponseWriter, r *http.Request) error {
	var id *Identification
	var err error
	id, err = identify(r)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	io.WriteString(w, id.JSON())
	return nil
}

func handleSets(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			return serve404(w)
		}
		vals, ok := r.Form["set"]
		if !ok {
			return serve404(w)
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		io.WriteString(w, "[\""+strings.Join(sets.Sets(vals...), "\", \"")+"\"]")
		return nil
	}
	k := sets.Keys()
	keys := make([]string, 0, len(k))
	for _, v := range k {
		if v[0] > 57 || v[0] < 48 {
			keys = append(keys, v)
		}
	}
	sort.Strings(keys)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return T("sets").Execute(w, keys)
}

func handleSiegfried(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return handleIdentify(w, r)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return T("siegfried").Execute(w, downloads)
}

func handleMain(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" || r.URL.Path != "/" {
		return serve404(w)
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return T("index").Execute(w, nil)
}
