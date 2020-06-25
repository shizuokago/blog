package datastore

import (
	"errors"
	"net/http"

	"cloud.google.com/go/datastore"
	"google.golang.org/appengine/user"
)

const KIND_USER = "User"

type User struct {
	Name      string
	Job       string
	Email     string
	URL       string
	TwitterId string
	AutoSave  bool
	Meta
}

func getUserKey(r *http.Request) *datastore.Key {
	c := r.Context()
	u := user.Current(c)

	return datastore.NameKey(KIND_USER, u.ID, nil)
}

func GetUser(r *http.Request) (*User, error) {

	c := r.Context()

	rtn := User{}
	key := getUserKey(r)

	client, err := createClient(c)

	err = client.Get(c, key, &rtn)
	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, err
		} else {
			return nil, nil
		}
	}
	return &rtn, nil
}

func SaveAvatar(r *http.Request) error {
	c := r.Context()
	u := user.Current(c)
	err := SaveFile(r, u.ID, FILE_TYPE_AVATAR)
	return err
}
