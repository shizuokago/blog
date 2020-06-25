package datastore

import (
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/pborman/uuid"
	"google.golang.org/api/iterator"
)

const KIND_ARTICLE = "Article"

type Article struct {
	Title       string
	SubTitle    string
	Tags        string
	PublishDate time.Time
	Markdown    []byte `datastore:",noindex"`
	Meta
}

func getArticleKey(r *http.Request, id string) *datastore.Key {
	return datastore.NameKey(KIND_ARTICLE, id, nil)
}

func SelectArticle(r *http.Request, p int) ([]Article, error) {

	c := r.Context()

	//TODO CURCOR

	q := datastore.NewQuery("Article").Order("- UpdatedAt").Limit(10)

	var s []Article

	client, err := createClient(c)
	if err != nil {
		return nil, err
	}

	t := client.Run(c, q)
	for {
		var a Article
		key, err := t.Next(&a)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		a.SetKey(key)
		s = append(s, a)
	}

	return s, nil
}

func GetArticle(r *http.Request, id string) (*Article, error) {
	c := r.Context()

	rtn := Article{}
	key := getArticleKey(r, id)

	client, err := createClient(c)
	rtn.SetKey(key)

	err = client.Get(c, key, &rtn)
	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &rtn, nil
}

func UpdateArticle(r *http.Request, id string, pub time.Time) (*Article, error) {

	r.ParseForm()
	title := r.FormValue("Title")
	tags := r.FormValue("Tags")
	mark := r.FormValue("Markdown")

	art, err := GetArticle(r, id)
	if err != nil {
		return nil, err
	}

	c := r.Context()

	art.Title = title
	art.SubTitle = CreateSubTitle(r.FormValue("Markdown"))
	art.Tags = tags
	art.Markdown = []byte(mark)
	if !pub.IsZero() {
		art.PublishDate = pub
	}

	client, err := createClient(c)
	_, err = client.Put(c, art.GetKey(), art)
	if err != nil {
		return nil, err
	}

	return art, nil
}

func DeleteArticle(r *http.Request, id string) error {

	c := r.Context()

	err := DeleteFile(r, "bg/"+id)
	if err != nil {
		return err
	}

	err = DeleteHtml(r, id)
	if err != nil {
		return err
	}

	akey := getArticleKey(r, id)

	client, err := createClient(c)
	err = client.Delete(c, akey)

	return err
}

func CreateArticle(r *http.Request) (string, error) {

	c := r.Context()
	id := uuid.New()

	bgd := GetBlog(r)
	base := bgd.Template
	article := &Article{
		Title:    "New Title",
		Tags:     bgd.Tags,
		Markdown: []byte(base),
	}

	article.Key = getArticleKey(r, id)

	client, err := createClient(c)
	_, err = client.Put(c, article.Key, article)
	if err != nil {
		return "", err
	}

	err = SaveFile(r, id, FILE_TYPE_BG)

	return id, err
}
