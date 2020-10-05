package datastore

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
)

var htmlCursor map[int]string

func init() {
	htmlCursor = make(map[int]string)
}

const KindHTML = "Html"

type HTMLParam struct {
	Body *HTML
	Data *HTMLData
}

type HTML struct {
	Title     string
	SubTitle  string
	Author    string
	AuthorID  string
	Updater   string
	UpdaterID string
	Meta
}

func getHTMLKey(key string) *datastore.Key {
	return datastore.NameKey(KindHTML, key, nil)
}

func GetHTML(ctx context.Context, k string) (*HTML, error) {

	rtn := HTML{}
	key := getHTMLKey(k)

	err := Get(ctx, key, &rtn)

	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		}
		return nil, xerrors.Errorf("datastore get: %w", err)
	}

	return &rtn, err
}

func UpdateHTML(ctx context.Context, key string, html *HTMLParam, art *Article) error {

	err := UpdateArticle(ctx, key, art)
	if err != nil {
		return xerrors.Errorf("update article: %w", err)
	}

	html.Body.Key = getHTMLKey(key)
	html.Data.SetKey(getHTMLDataKey(key))

	err = Put(ctx, html.Body)
	if err != nil {
		return xerrors.Errorf("put html: %w", err)
	}

	err = Put(ctx, html.Data)
	if err != nil {
		return xerrors.Errorf("put html data: %w", err)
	}

	//再度更新
	htmlCursor = make(map[int]string)
	return nil
}

func SelectHTML(ctx context.Context, p int) ([]HTML, error) {

	q := datastore.NewQuery(KindHTML).Order("- CreatedAt").Limit(5)

	var s []HTML

	client, err := createClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("get html: %w", err)
	}

	if cur, ok := htmlCursor[p]; ok {
		cursor, err := datastore.DecodeCursor(cur)
		if err != nil {
			log.Printf("%+v", err)
		} else {
			q = q.Start(cursor)
		}
	}

	t := client.Run(ctx, q)
	for {
		var h HTML
		key, err := t.Next(&h)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf("html(next): %w", err)
		}
		h.SetKey(key)
		s = append(s, h)
	}

	next, err := t.Cursor()
	if err != nil {
		log.Printf("%+v", err)
	} else {
		htmlCursor[p+1] = next.String()
	}

	return s, nil
}

func DeleteHTML(ctx context.Context, id string) error {

	hkey := getHTMLKey(id)

	client, err := createClient(ctx)

	err = client.Delete(ctx, hkey)
	if err != nil {
		return xerrors.Errorf("delete html: %w", err)
	}
	hdkey := getHTMLDataKey(id)
	err = client.Delete(ctx, hdkey)
	if err != nil {
		return xerrors.Errorf("delete html data: %w", err)
	}

	htmlCursor = make(map[int]string)
	return nil
}

const KindHTMLData = "HtmlData"

type HTMLData struct {
	key     *datastore.Key
	Content []byte `datastore:",noindex"`
}

func getHTMLDataKey(key string) *datastore.Key {
	return datastore.NameKey(KindHTMLData, key, nil)
}

func (d *HTMLData) GetKey() *datastore.Key {
	return d.key
}

func (d *HTMLData) SetKey(k *datastore.Key) {
	d.key = k
}

func GetHTMLData(ctx context.Context, k string) (*HTMLData, error) {

	rtn := HTMLData{}
	key := getHTMLDataKey(k)

	client, err := createClient(ctx)
	err = client.Get(ctx, key, &rtn)

	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		} else {
			return nil, xerrors.Errorf("get html data: %w", err)
		}
	}
	return &rtn, nil
}
