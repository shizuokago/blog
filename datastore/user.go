package datastore

import (
	"context"
	"errors"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
)

const KindUser = "User"

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
	return datastore.NameKey(KindUser, key, nil)
}

func GetUser(ctx context.Context, email string) (*User, error) {

	rtn := User{}
	key := getUserKey(email)

	err := Get(ctx, key, &rtn)
	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		} else {
			return nil, xerrors.Errorf("get user: %w", err)
		}
	}
	return &rtn, nil
}

func SaveAvatar(ctx context.Context, key string, p *FileParam) error {

	err := SaveFile(ctx, key, FILE_TYPE_AVATAR, p)
	if err != nil {
		return xerrors.Errorf("save file: %w", err)
	}
	return nil
}
