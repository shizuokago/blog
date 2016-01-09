package main

import (
	verr "github.com/knightso/base/errors"
	"github.com/knightso/base/gae/ds"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

func init() {
	ds.DefaultCache = true
}

const KIND_USER = "User"

type User struct {
	Name      string
	Job       string
	Email     string
	URL       string
	TwitterId string
	ds.Meta
}

func getUserKey(r *http.Request) *datastore.Key {
	c := appengine.NewContext(r)
	u := user.Current(c)
	return datastore.NewKey(c, KIND_USER, u.ID, 0, nil)
}

func getUser(r *http.Request) (*User, error) {

	c := appengine.NewContext(r)

	rtn := User{}
	rtn.Key = getUserKey(r)

	err := ds.Get(c, rtn.Key, &rtn)
	if err != nil && verr.Root(err) != datastore.ErrNoSuchEntity {
		return nil, verr.Root(err)
	}
	return &rtn, nil
}

func putUser(r *http.Request) (*User, error) {

	c := appengine.NewContext(r)

	r.ParseForm()
	rtn := User{
		Name:      r.FormValue("Name"),
		Job:       r.FormValue("Job"),
		Email:     r.FormValue("Email"),
		URL:       r.FormValue("Url"),
		TwitterId: r.FormValue("TwitterId"),
	}

	rtn.Key = getUserKey(r)
	err := ds.Put(c, &rtn)
	if err != nil {
		return nil, err
	}
	return &rtn, nil
}

const KIND_ARTICLE = "Article"

type Article struct {
	Title    string
	SubTitle string
	Tags     string
	Markdown string
	ds.Meta
}

const KIND_HTML = "Html"

type Html struct {
	Content string
	ds.Meta
}

const KIND_FILE = "File"

type File struct {
	Name string
	Size int64
	Mime string
	ds.Meta
}

const KIND_FILEDATA = "FileData"

type FileData struct {
	Content []byte
}
