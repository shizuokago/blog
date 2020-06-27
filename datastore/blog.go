package datastore

import (
	"log"
	"net/http"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
)

const KIND_BLOG = "Blog"

type Blog struct {
	Name        string
	Author      string
	Tags        string
	Description string
	Template    string
	Meta
}

var pkgBlog = Blog{}

func GetBlog(r *http.Request) *Blog {

	if pkgBlog.Name != "" {
		return &pkgBlog
	}
	c := r.Context()
	key := datastore.NameKey(KIND_BLOG, "Fixing", nil)

	err := Get(c, key, &pkgBlog)
	if err != nil {
		log.Println(err)
	}
	return &pkgBlog
}

func PutBlog(r *http.Request) error {

	pkgBlog = Blog{
		Name:        r.FormValue("BlogName"),
		Author:      r.FormValue("BlogAuthor"),
		Description: r.FormValue("Description"),
		Tags:        r.FormValue("BlogTags"),
		Template:    r.FormValue("BlogTemplate"),
	}

	c := r.Context()
	key := datastore.NameKey(KIND_BLOG, "Fixing", nil)

	pkgBlog.SetKey(key)

	err := Put(c, &pkgBlog)
	if err != nil {
		return xerrors.Errorf("blog put: %w", err)
	}

	return nil
}
