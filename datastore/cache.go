package datastore

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/pborman/uuid"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
)

type Cache struct {
	Pattern string
	MaxAge  int
	Meta
}

const (
	KindCache = "Cache"
)

func createCacheKey() *datastore.Key {
	id := uuid.New()
	return datastore.NameKey(KindCache, id, nil)
}

func GetCacheKey(name string) *datastore.Key {
	return datastore.NameKey(KindCache, name, nil)
}

func initCache(ctx context.Context) error {

	_, err := RegisterCache(ctx, "/js/*", 3600)
	if err != nil {
		return xerrors.Errorf("/js/* cache register: %w", err)
	}

	_, err = RegisterCache(ctx, "/css/*", 3600)
	if err != nil {
		return xerrors.Errorf("/css/* cache register: %w", err)
	}

	_, err = RegisterCache(ctx, "/static/css/*", 3600)
	if err != nil {
		return xerrors.Errorf("/static/css/*(deprecated) cache register: %w", err)
	}

	_, err = RegisterCache(ctx, "/images/*", 3600)
	if err != nil {
		return xerrors.Errorf("/images/* cache register: %w", err)
	}

	_, err = RegisterCache(ctx, "/file/*/*", 3600)
	if err != nil {
		return xerrors.Errorf("/file/* cache register: %w", err)
	}

	_, err = RegisterCache(ctx, "/entry/*", 3600)
	if err != nil {
		return xerrors.Errorf("/entry/* cache register: %w", err)
	}

	_, err = RegisterCache(ctx, "/favicon.ico", 3600)
	if err != nil {
		return xerrors.Errorf("/favicon.ico cache register: %w", err)
	}

	_, err = RegisterCache(ctx, "/", 600)
	if err != nil {
		return xerrors.Errorf("/ cache register: %w", err)
	}

	return nil
}

func RegisterCache(ctx context.Context, p string, age int) (*Cache, error) {

	cache := Cache{}

	id := createCacheKey()

	cache.SetKey(id)
	cache.Pattern = p
	cache.MaxAge = age

	err := Put(ctx, &cache)
	if err != nil {
		return nil, xerrors.Errorf("put article: %w", err)
	}

	return &cache, nil
}

func SelectCache(ctx context.Context) ([]*Cache, error) {

	var s []*Cache

	q := datastore.NewQuery(KindCache)

	client, err := createClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("create client: %w", err)
	}

	t := client.Run(ctx, q)
	for {
		var c Cache
		key, err := t.Next(&c)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf("next: %w", err)
		}
		c.SetKey(key)
		s = append(s, &c)
	}

	return s, nil
}

func GetCaches(ctx context.Context) ([]*Cache, error) {

	cache, err := SelectCache(ctx)
	if err != nil {
		return nil, xerrors.Errorf("select cache: %w", err)
	}

	if len(cache) == 0 {
		err = initCache(ctx)
		if err != nil {
			return nil, xerrors.Errorf("initCache: %w", err)
		}
		return SelectCache(ctx)
	}

	return cache, nil
}
