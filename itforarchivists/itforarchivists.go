package itforarchivists

import (
	"net/http"
	"os"
	"time"

	"github.com/richardlehane/siegfried"
)

func init() {
	// setup gobsize
	f, _ := os.Open("public/latest/default.sig")
	i, _ := f.Stat()
	current.Size = int(i.Size())
	f.Close()

	// setup created
	s, _ := siegfried.Load("public/latest/default.sig")
	current.Created = s.C.Format(time.RFC3339)

	// routes
	http.HandleFunc("/", wrapError(handleMain))
	http.HandleFunc("/siegfried", wrapError(handleSiegfried))
	http.HandleFunc("/siegfried/update", handleUpdate)
	http.HandleFunc("/siegfried/sets", wrapError(handleSets))

	// LATEST SIG
	// last modified header
	gmt, _ := time.LoadLocation("GMT")
	lastModified := s.C.In(gmt)
	gobbler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		f, _ := os.Open("public/latest/default.sig")
		// todo - get rid of lastModified, file opening etc. and replace with http.ServeFile when the AppEngine lastmodified works
		http.ServeContent(w, r, "default.sig", lastModified, f)
		f.Close()
	}
	http.HandleFunc("/siegfried/latest", gobbler)
}
