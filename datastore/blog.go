package datastore

import (
	"net/http"

	"cloud.google.com/go/datastore"
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

	client, err := createClient(c)
	if err != nil {
		//
	}

	err = client.Get(c, key, &pkgBlog)
	if err != nil {
		// Nothing
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

	client, err := createClient(c)

	_, err = client.Put(c, key, &pkgBlog)
	if err != nil {
		return err
	}

	return nil
}
