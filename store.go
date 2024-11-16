package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

const delimiter = "/"

type store interface {
	stash(uuid, name, title, desc string, body interface{}) error
	retrieve(uuid string) (name, title, desc string, body []byte, err error)
	list(prefix, dir string) []string
	dirs(prefix string) []string
}

func newStore(r *http.Request) (store, error) {
	if os.Getenv("GAE_APPLICATION") == "" { // detect if running locally
		return newSimpleStore(r)
	}
	return newCloudStore(r)
}

var sStore simpleStore

func newSimpleStore(r *http.Request) (simpleStore, error) {
	if sStore == nil {
		sStore = make(simpleStore)
	}
	return sStore, nil
}

type simpleStore map[string]struct {
	name  string
	title string
	desc  string
	body  []byte
}

func (ss simpleStore) stash(key, name, title, desc string, body interface{}) error {
	byt, err := json.Marshal(body)
	if err != nil {
		return err
	}
	if _, ok := ss[key]; ok {
		return fmt.Errorf("attempted to store same key twice: %s", key)
	}
	ss[key] = struct {
		name  string
		title string
		desc  string
		body  []byte
	}{
		name:  name,
		title: title,
		desc:  desc,
		body:  byt,
	}
	return nil
}

func (ss simpleStore) retrieve(key string) (string, string, string, []byte, error) {
	v, ok := ss[key]
	if !ok {
		return "", "", "", nil, fmt.Errorf("can't find key: %s", key)
	}
	return v.name, v.title, v.desc, v.body, nil
}

func (ss simpleStore) list(prefix, dir string) []string {
	path := prefix + delimiter + dir + delimiter
	ret := make([]string, 0, 100)
	for k := range ss {
		if strings.HasPrefix(k, path) {
			ret = append(ret, k)
		}
	}
	return ret
}

func (ss simpleStore) dirs(prefix string) []string {
	uniques := make(map[string]struct{})
	for k := range ss {
		if strings.HasPrefix(k, prefix+delimiter) && strings.Count(k, delimiter) >= 2 {
			uniques[k[len(prefix+delimiter):strings.Index(k[len(prefix+delimiter):], delimiter)+len(prefix+delimiter)]] = struct{}{}
		}
	}
	ret := make([]string, 0, len(uniques))
	for k := range uniques {
		ret = append(ret, k)
	}
	return ret
}

type cloudStore struct {
	bucket string
	ctx    context.Context
}

func newCloudStore(r *http.Request) (*cloudStore, error) {
	return &cloudStore{
		bucket: "itforarchivists.appspot.com", // hardcode this in as can't get the default anymore without appengine pkg
		ctx:    r.Context(),
	}, nil
}

func (cs *cloudStore) stash(key, name, title, desc string, body interface{}) error {
	byt, err := json.Marshal(body)
	if err != nil {
		return err
	}
	client, err := storage.NewClient(cs.ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	obj := client.Bucket(cs.bucket).Object(key)
	if rc, err := obj.NewReader(cs.ctx); err != storage.ErrObjectNotExist {
		if err == nil {
			rc.Close()
			err = fmt.Errorf("can't store same key twice: %s", key)
		}
		return err
	}
	wc := obj.NewWriter(cs.ctx)
	wc.ContentType = "application/json"
	wc.Metadata = map[string]string{
		"x-goog-meta-name":        name,
		"x-goog-meta-title":       title,
		"x-goog-meta-description": desc,
	}
	if _, err := wc.Write(byt); err != nil {
		wc.Close()
		return err
	}
	return wc.Close()
}

func (cs *cloudStore) retrieve(key string) (string, string, string, []byte, error) {
	client, err := storage.NewClient(cs.ctx)
	if err != nil {
		return "", "", "", nil, err
	}
	defer client.Close()
	obj := client.Bucket(cs.bucket).Object(key)
	rc, err := obj.NewReader(cs.ctx)
	if err != nil {
		return "", "", "", nil, err
	}
	defer rc.Close()
	body, err := io.ReadAll(rc)
	if err != nil {
		return "", "", "", nil, err
	}
	var name, title, desc string
	attrs, err := obj.Attrs(cs.ctx)
	if err == nil && attrs.Metadata != nil {
		name, title, desc = attrs.Metadata["x-goog-meta-name"], attrs.Metadata["x-goog-meta-title"], attrs.Metadata["x-goog-meta-description"]
	}
	return name, title, desc, body, nil
}

func (cs *cloudStore) list(prefix, dir string) []string {
	client, err := storage.NewClient(cs.ctx)
	if err != nil {
		return nil
	}
	defer client.Close()
	path := prefix + delimiter + dir + delimiter
	query := &storage.Query{Prefix: path}
	ret := make([]string, 0, 20)
	it := client.Bucket(cs.bucket).Objects(cs.ctx, query)
	for obj, err := it.Next(); err == nil; obj, err = it.Next() {
		ret = append(ret, obj.Name)
	}
	return ret
}

func (cs *cloudStore) dirs(prefix string) []string {
	client, err := storage.NewClient(cs.ctx)
	if err != nil {
		return nil
	}
	defer client.Close()
	query := &storage.Query{
		Prefix: prefix + delimiter,
	}
	uniques := make(map[string]struct{})
	it := client.Bucket(cs.bucket).Objects(cs.ctx, query)
	for obj, err := it.Next(); err == nil; obj, err = it.Next() {
		k := obj.Name
		if strings.HasPrefix(k, prefix+delimiter) && strings.Count(k, delimiter) >= 2 {
			uniques[k[len(prefix+delimiter):strings.Index(k[len(prefix+delimiter):], delimiter)+len(prefix+delimiter)]] = struct{}{}
		}
	}
	ret := make([]string, 0, len(uniques))
	for k := range uniques {
		ret = append(ret, k)
	}
	return ret
}
