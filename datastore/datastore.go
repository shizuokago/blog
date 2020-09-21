package datastore

import (
	"context"
	"strings"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"

	"github.com/shizuokago/blog/config"
)

func init() {
}

func createClient(ctx context.Context) (*datastore.Client, error) {

	client, err := datastore.NewClient(ctx, config.ProjectID())
	if err != nil {
		return nil, xerrors.Errorf("datastore client: %w", err)
	}
	return client, nil
}

func CreateSubTitle(src string) string {

	dst := strings.Replace(src, "\n", "", -1)
	dst = strings.Replace(dst, "*", "", -1)

	if len(dst) > 600 {
		dst = string([]rune(dst)[0:200]) + "..."
	}
	return dst
}

func Get(ctx context.Context, key *datastore.Key, dst HasKey) error {

	client, err := createClient(ctx)
	if err != nil {
		return xerrors.Errorf("create client: %w", err)
	}

	err = client.Get(ctx, key, dst)
	if err != nil {
		return xerrors.Errorf("datastore get: %w", err)
	}

	dst.SetKey(key)

	return nil
}

func Put(ctx context.Context, dst HasKey) error {

	if t, ok := dst.(HasTime); ok {
		t.SetTime()
	}
	if v, ok := dst.(HasVersion); ok {
		v.IncrementVersion()
	}

	key := dst.GetKey()
	if key == nil {
		return xerrors.Errorf("datastore put error -> key is nil")
	}

	client, err := createClient(ctx)
	if err != nil {
		return xerrors.Errorf("create client: %w", err)
	}

	_, err = client.Put(ctx, key, dst)
	if err != nil {
		return xerrors.Errorf("datastore put error: %w", err)
	}
	return nil
}

func Delete(ctx context.Context, key *datastore.Key) error {

	client, err := createClient(ctx)
	if err != nil {
		return xerrors.Errorf("create client: %w", err)
	}

	err = client.Delete(ctx, key)
	if err != nil {
		return xerrors.Errorf("datastore delete error: %w", err)
	}
	return nil
}
