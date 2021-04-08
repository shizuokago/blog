package datastore

import (
	"context"
	"log"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
)

const (
	BlogKey  = "Fixing"
	KindBlog = "Blog"
)

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

func GetBlog(ctx context.Context) *Blog {

	if pkgBlog.Name != "" {
		return &pkgBlog
	}

	key := datastore.NameKey(KindBlog, BlogKey, nil)

	err := Get(ctx, key, &pkgBlog)
	if err != nil {
		log.Println(err)
	}
	return &pkgBlog
}

func PutBlog(ctx context.Context, blog *Blog) error {

	key := datastore.NameKey(KindBlog, BlogKey, nil)
	blog.SetKey(key)

	err := Put(ctx, blog)
	if err != nil {
		return xerrors.Errorf("blog put: %w", err)
	}

	pkgBlog = *blog

	return nil
}

func (b Blog) getUsers() []string {

	if b.Users == "" {
		return nil
	}
	users := strings.Split(b.Users, ",")
	return users
}

func IsUser(ctx context.Context, id string) bool {

	b := GetBlog(ctx)

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

func PutInformation(ctx context.Context, blog *Blog, user *User, mail string) error {

	blog.SetKey(datastore.NameKey(KindBlog, BlogKey, nil))
	err := Put(ctx, blog)
	if err != nil {
		return xerrors.Errorf("put blog: %w", err)
	}

	user.SetKey(getUserKey(mail))
	err = Put(ctx, user)
	if err != nil {
		return xerrors.Errorf("put user: %w", err)
	}
	return nil
}
