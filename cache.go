package main

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"
)

type cachedPage struct {
	t      time.Time
	status int
	header http.Header
	buf    *bytes.Buffer
}

func (cp *cachedPage) Header() http.Header {
	return cp.header
}

func (cp *cachedPage) Write(b []byte) (int, error) {
	return cp.buf.Write(b)
}

func (cp *cachedPage) WriteHeader(statusCode int) {
	cp.status = statusCode
}

type cache struct {
	du   time.Duration
	cmap map[string]*cachedPage
	mu   *sync.Mutex
}

func newCache(du time.Duration) *cache {
	return &cache{du, make(map[string]*cachedPage), &sync.Mutex{}}
}

func hdlCache(f func(http.ResponseWriter, *http.Request), c *cache) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c.mu.Lock()
		defer c.mu.Unlock()
		cp, ok := c.cmap[r.URL.Path]
		if !ok || time.Since(cp.t) > c.du {
			cp = &cachedPage{time.Now(), 0, http.Header{}, &bytes.Buffer{}}
			f(cp, r)
			c.cmap[r.URL.Path] = cp
		}
		for k, v := range cp.header {
			w.Header()[k] = v
		}
		if cp.status > 0 {
			w.WriteHeader(cp.status)
		}
		io.Copy(w, bytes.NewReader(cp.buf.Bytes()))
	}
}
