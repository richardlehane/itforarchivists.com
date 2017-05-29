package itforarchivists

import (
	"html/template"
	"path/filepath"
	"sync"
)

var cachedTemplates = map[string]*template.Template{}
var cachedMutex sync.Mutex

func T(name string) *template.Template {
	cachedMutex.Lock()
	defer cachedMutex.Unlock()

	if t, ok := cachedTemplates[name]; ok {
		return t
	}
	t := template.Must(template.ParseFiles(
		filepath.Join("templates", "base.html"),
		filepath.Join("templates", name+".html"),
	))
	cachedTemplates[name] = t

	return t
}
