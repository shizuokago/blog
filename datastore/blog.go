package datastore

import (
	"log"
	"net/http"
	"strings"

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
	Users       string
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
		Users:       r.FormValue("Users"),
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

func (b Blog) getUsers() []string {

	if b.Users == "" {
		return nil
	}
	users := strings.Split(b.Users, ",")
	return users
}

func IsUser(r *http.Request, id string) bool {

	b := GetBlog(r)

	users := b.getUsers()
	if users == nil {
		log.Println("ユーザが存在しない為、フルアクセス")
		return true
	}

	//ユーザが存在した場合OK
	for _, elm := range users {
		if elm == id {
			return true
		}
	}

	return false
}
