package blog

import (
	verr "github.com/knightso/base/errors"
	"github.com/knightso/base/gae/ds"
	"io/ioutil"
	"net/http"

	"github.com/pborman/uuid"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/user"
)

func init() {
	ds.DefaultCache = true
}

const KIND_BLOG = "Blog"

type Blog struct {
	Name string
	ds.Meta
}

const KIND_USER = "User"

type User struct {
	Name      string
	Job       string
	Email     string
	URL       string
	TwitterId string
	ds.Meta
}

func getUserKey(r *http.Request) *datastore.Key {
	c := appengine.NewContext(r)
	u := user.Current(c)
	return datastore.NewKey(c, KIND_USER, u.ID, 0, nil)
}

func getUser(r *http.Request) (*User, error) {

	c := appengine.NewContext(r)

	rtn := User{}
	key := getUserKey(r)

	err := ds.Get(c, key, &rtn)
	if err != nil && verr.Root(err) != datastore.ErrNoSuchEntity {
		return nil, verr.Root(err)
	}
	return &rtn, nil
}

func putUser(r *http.Request) (*User, error) {

	c := appengine.NewContext(r)

	r.ParseForm()
	rtn := User{
		Name:      r.FormValue("Name"),
		Job:       r.FormValue("Job"),
		Email:     r.FormValue("Email"),
		URL:       r.FormValue("Url"),
		TwitterId: r.FormValue("TwitterId"),
	}

	rtn.Key = getUserKey(r)
	err := ds.Put(c, &rtn)
	if err != nil {
		return nil, err
	}
	return &rtn, nil
}

const KIND_ARTICLE = "Article"

type Article struct {
	Title    string
	SubTitle string
	Tags     string
	Markdown datastore.ByteString
	ds.Meta
}

func getArticleKey(r *http.Request, id string) *datastore.Key {
	c := appengine.NewContext(r)
	return datastore.NewKey(c, KIND_ARTICLE, id, 0, nil)
}

func selectArticle(r *http.Request, page int) ([]Article, error) {

	c := appengine.NewContext(r)

	q := datastore.NewQuery("Article").
		Order("- UpdatedAt")

	var s []Article
	err := ds.ExecuteQuery(c, q, &s)

	//TODO 違う
	if err != nil && verr.Root(err) != datastore.ErrNoSuchEntity {
		return nil, verr.Root(err)
	}

	return s, nil
}

func getArticle(r *http.Request, id string) (*Article, error) {
	c := appengine.NewContext(r)

	rtn := Article{}
	key := getArticleKey(r, id)

	err := ds.Get(c, key, &rtn)
	if err != nil && verr.Root(err) != datastore.ErrNoSuchEntity {
		return nil, verr.Root(err)
	}
	return &rtn, nil
}

const KIND_HTML = "Html"

type Html struct {
	Content datastore.ByteString
	ds.Meta
}

const KIND_FILE = "File"

type File struct {
	Size int64
	ds.Meta
}

func getFileKey(r *http.Request, name string) *datastore.Key {
	c := appengine.NewContext(r)
	return datastore.NewKey(c, KIND_FILE, name, 0, nil)
}

func createArticle(r *http.Request) (string, error) {

	upload, header, err := r.FormFile("file")
	if err != nil {
		//add error handling
		return "", err
	}
	defer upload.Close()

	b, err := ioutil.ReadAll(upload)
	if err != nil {
		return "", err
	}

	c := appengine.NewContext(r)
	id := uuid.New()

	article := &Article{
		Title:    "New Title",
		Markdown: []byte("* Section1"),
	}
	article.Key = getArticleKey(r, id)
	err = ds.Put(c, article)
	if err != nil {
		return "", err
	}

	file := &File{
		Size: int64(len(b)),
	}

	file.Key = getFileKey(r, id)
	err = ds.Put(c, file)
	if err != nil {
		return "", err
	}

	fileData := &FileData{
		Content: b,
		Mime:    header.Header["Content-Type"][0],
	}
	fileData.SetKey(getFileDataKey(r, id))
	err = ds.Put(c, fileData)
	if err != nil {
		return "", err
	}

	return id, nil
}

const KIND_FILEDATA = "FileData"

type FileData struct {
	key     *datastore.Key
	Mime    string
	Content []byte
}

func (d *FileData) GetKey() *datastore.Key {
	return d.key
}

func (d *FileData) SetKey(k *datastore.Key) {
	d.key = k
}

func getFileDataKey(r *http.Request, name string) *datastore.Key {
	c := appengine.NewContext(r)
	return datastore.NewKey(c, KIND_FILEDATA, name, 0, nil)
}

func getFileData(r *http.Request, name string) (*FileData, error) {
	c := appengine.NewContext(r)

	rtn := FileData{}
	key := getFileDataKey(r, name)

	err := ds.Get(c, key, &rtn)
	if err != nil && verr.Root(err) != datastore.ErrNoSuchEntity {
		return nil, verr.Root(err)
	}
	return &rtn, nil
}