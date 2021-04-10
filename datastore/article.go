package datastore

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/pborman/uuid"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
)

var articleCursor map[int]string

func init() {
	articleCursor = make(map[int]string)
}

const KindArticle = "Article"

type Article struct {
	Title       string
	SubTitle    string
	Tags        string
	PublishDate time.Time
	Markdown    []byte `datastore:",noindex"`
	Meta
}

func getArticleKey(id string) *datastore.Key {
	return datastore.NameKey(KindArticle, id, nil)
}

func SelectArticle(ctx context.Context, p int) ([]Article, error) {

	q := datastore.NewQuery(KindArticle).Order("- UpdatedAt").Limit(10)

	var s []Article

	client, err := createClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("create client: %w", err)
	}

	if cur, ok := articleCursor[p]; ok {
		cursor, err := datastore.DecodeCursor(cur)
		if err != nil {
			log.Printf("%+v", err)
		} else {
			q = q.Start(cursor)
		}
	}

	t := client.Run(ctx, q)
	for {
		var a Article
		key, err := t.Next(&a)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf("next: %w", err)
		}
		a.SetKey(key)
		s = append(s, a)
	}

	next, err := t.Cursor()
	if err != nil {
		log.Printf("%+v", err)
	} else {
		articleCursor[p+1] = next.String()
	}

	return s, nil
}

func GetArticle(ctx context.Context, id string) (*Article, error) {

	rtn := Article{}
	key := getArticleKey(id)

	err := Get(ctx, key, &rtn)
	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		} else {
			return nil, xerrors.Errorf("client get: %w", err)
		}
	}

	return &rtn, nil
}

func UpdateArticle(ctx context.Context, id string, art *Article) error {

	key := getArticleKey(id)
	art.SetKey(key)

	err := Put(ctx, art)
	if err != nil {
		return xerrors.Errorf("put article: %w", err)
	}

	return nil
}

func DeleteArticle(ctx context.Context, id string) error {

	err := DeleteFile(ctx, "bg/"+id)
	if err != nil {
		return err
	}

	err = DeleteHTML(ctx, id)
	if err != nil {
		return err
	}

	akey := getArticleKey(id)

	client, err := createClient(ctx)
	err = client.Delete(ctx, akey)
	if err != nil {
		return xerrors.Errorf("delete article: %w", err)
	}

	articleCursor = make(map[int]string)
	return nil
}

func CreateArticle(ctx context.Context, art *Article, f *File, d *FileData) (string, error) {

	id := uuid.New()
	art.Key = getArticleKey(id)

	err := Put(ctx, art)
	if err != nil {
		return "", xerrors.Errorf("Article put: %w", err)
	}

	p := FileParam{
		File:     f,
		FileData: d,
	}

	err = SaveFile(ctx, id, FileTypeBG, &p)
	if err != nil {
		return "", xerrors.Errorf("File put: %w", err)
	}

	articleCursor = make(map[int]string)

	return id, nil
}
