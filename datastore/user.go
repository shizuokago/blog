package datastore

import (
	"errors"
	"net/http"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
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

	err := Get(c, key, &rtn)
	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		} else {
			return nil, xerrors.Errorf("get user: %w", err)
		}
	}
	return &rtn, nil
}

func SaveAvatar(r *http.Request, key string) error {
	err := SaveFile(r, key, FILE_TYPE_AVATAR)
	if err != nil {
		return xerrors.Errorf("save file: %w", err)
	}
	return nil
}
