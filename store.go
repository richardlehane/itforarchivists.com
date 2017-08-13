package itforarchivists

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/file"
)

type store interface {
	stash(uuid, name, title, desc string, res *Results) error
	retrieve(uuid string) (name, title, desc string, res []byte, err error)
}

type simpleStore map[string]struct {
	name  string
	title string
	desc  string
	res   []byte
}

func (ss simpleStore) stash(uuid, name, title, desc string, res *Results) error {
	byt, err := json.Marshal(res)
	if err != nil {
		return err
	}
	if _, ok := ss[uuid]; ok {
		return fmt.Errorf("attempted to store same uuid twice: %s", uuid)
	}
	ss[uuid] = struct {
		name  string
		title string
		desc  string
		res   []byte
	}{
		name:  name,
		title: title,
		desc:  desc,
		res:   byt,
	}
	return nil
}

func (ss simpleStore) retrieve(uuid string) (string, string, string, []byte, error) {
	v, ok := ss[uuid]
	if !ok {
		return "", "", "", nil, fmt.Errorf("can't find results: %s", uuid)
	}
	return v.name, v.title, v.desc, v.res, nil
}

type cloudStore struct {
	bucket string
	ctx    context.Context
}

func newCloudStore(r *http.Request) (*cloudStore, error) {
	ctx := appengine.NewContext(r)
	bucket, err := file.DefaultBucketName(ctx)
	if err != nil {
		return nil, err
	}
	return &cloudStore{
		bucket: bucket,
		ctx:    ctx,
	}, nil
}

func (cs *cloudStore) stash(uuid, name, title, desc string, res *Results) error {
	byt, err := json.Marshal(res)
	if err != nil {
		return err
	}
	client, err := storage.NewClient(cs.ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	obj := client.Bucket(cs.bucket).Object(uuid)
	if rc, err := obj.NewReader(cs.ctx); err != storage.ErrObjectNotExist {
		if err == nil {
			rc.Close()
			err = fmt.Errorf("can't store same uuid twice: %s", uuid)
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

func (cs *cloudStore) retrieve(uuid string) (string, string, string, []byte, error) {
	client, err := storage.NewClient(cs.ctx)
	if err != nil {
		return "", "", "", nil, err
	}
	defer client.Close()
	obj := client.Bucket(cs.bucket).Object(uuid)
	rc, err := obj.NewReader(cs.ctx)
	if err != nil {
		return "", "", "", nil, err
	}
	defer rc.Close()
	res, err := ioutil.ReadAll(rc)
	if err != nil {
		return "", "", "", nil, err
	}
	var name, title, desc string
	attrs, err := obj.Attrs(cs.ctx)
	if err == nil && attrs.Metadata != nil {
		name, title, desc = attrs.Metadata["x-goog-meta-name"], attrs.Metadata["x-goog-meta-title"], attrs.Metadata["x-goog-meta-description"]
	}
	return name, title, desc, res, nil
}
