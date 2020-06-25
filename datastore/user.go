package datastore

import (
	"errors"
	"net/http"

	"cloud.google.com/go/datastore"
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

func getUserKey(key string) *datastore.Key {
	return datastore.NameKey(KIND_USER, key, nil)
}

func GetUser(r *http.Request, email string) (*User, error) {

	c := r.Context()

	rtn := User{}
	key := getUserKey(email)

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

func SaveAvatar(r *http.Request, key string) error {
	err := SaveFile(r, key, FILE_TYPE_AVATAR)
	return err
}
