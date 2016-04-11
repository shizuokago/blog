package blog

import (
	"net/http"
	"strings"
	"time"

	verr "github.com/knightso/base/errors"
	"github.com/knightso/base/gae/ds"
	"github.com/pborman/uuid"

	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
)

func init() {
	ds.DefaultCache = true
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
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return nil, verr.Root(err)
		} else {
			return nil, nil
		}
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
	Title       string
	SubTitle    string
	Tags        string
	PublishDate time.Time
	Markdown    datastore.ByteString `datastore:",noindex"`
	ds.Meta
}

func getArticleKey(r *http.Request, id string) *datastore.Key {
	c := appengine.NewContext(r)
	return datastore.NewKey(c, KIND_ARTICLE, id, 0, nil)
}

func selectArticle(r *http.Request, p int) ([]Article, error) {

	c := appengine.NewContext(r)
	item, err := memcache.Get(c, "article_"+strconv.Itoa(p)+"_cursor")
	cursor := ""
	if err == nil {
		cursor = string(item.Value)
	}

	q := datastore.NewQuery("Article").
		Order("- UpdatedAt").
		Limit(10)

	if cursor != "" {
		cur, err := datastore.DecodeCursor(cursor)
		if err == nil {
			q = q.Start(cur)
		}
	}

	var s []Article
	t := q.Run(c)
	for {
		var a Article
		key, err := t.Next(&a)

		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		a.SetKey(key)
		s = append(s, a)
	}

	cur, err := t.Cursor()
	if err != nil {
		return nil, err
	}

	err = memcache.Set(c, &memcache.Item{
		Key:   "article_" + strconv.Itoa(p+1) + "_cursor",
		Value: []byte(cur.String()),
	})

	if err != nil {
		return nil, err
	}

	return s, nil
}

func getArticle(r *http.Request, id string) (*Article, error) {
	c := appengine.NewContext(r)

	rtn := Article{}
	key := getArticleKey(r, id)

	err := ds.Get(c, key, &rtn)
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return nil, verr.Root(err)
		} else if verr.Root(err) == datastore.ErrNoSuchEntity {
			return nil, nil
		}
	}
	return &rtn, nil
}

func updateArticle(r *http.Request, id string, pub time.Time) (*Article, error) {

	r.ParseForm()
	title := r.FormValue("Title")
	tags := r.FormValue("Tags")
	mark := datastore.ByteString(r.FormValue("Markdown"))

	art, err := getArticle(r, id)
	if err != nil {
		return nil, err
	}

	c := appengine.NewContext(r)

	art.Title = title
	art.SubTitle = createSubTitle(r.FormValue("Markdown"))
	art.Tags = tags
	art.Markdown = mark
	if !pub.IsZero() {
		art.PublishDate = pub
	}

	err = ds.Put(c, art)
	if err != nil {
		return nil, err
	}

	return art, nil
}

func createSubTitle(src string) string {

	dst := strings.Replace(src, "\n", "", -1)
	dst = strings.Replace(dst, "*", "", -1)

	if len(dst) > 600 {
		dst = string([]rune(dst)[0:200]) + "..."
	}
	return dst
}

func createArticle(r *http.Request) (string, error) {

	c := appengine.NewContext(r)
	id := uuid.New()

	base := blog.Template
	article := &Article{
		Title:    "New Title",
		Tags:     blog.Tags,
		Markdown: []byte(base),
	}

	article.Key = getArticleKey(r, id)
	err := ds.Put(c, article)
	if err != nil {
		return "", err
	}

	err = saveFile(r, id, FILE_TYPE_BG)

	return id, err
}

func deleteArticle(r *http.Request, id string) error {

	c := appengine.NewContext(r)

	err := deleteFile(r, "bg/"+id)
	if err != nil {
		return err
	}

	err = deleteHtml(r, id)
	if err != nil {
		return err
	}

	akey := getArticleKey(r, id)
	err = ds.Delete(c, akey)

	return err
}

const KIND_HTML = "Html"

type Html struct {
	Title    string
	SubTitle string
	Author   string
	AuthorID string
	ds.Meta
}

func getHtmlKey(r *http.Request, key string) *datastore.Key {
	c := appengine.NewContext(r)
	return datastore.NewKey(c, KIND_HTML, key, 0, nil)
}

func getHtml(r *http.Request, k string) (*Html, error) {

	c := appengine.NewContext(r)

	rtn := Html{}
	key := getHtmlKey(r, k)
	err := ds.Get(c, key, &rtn)
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return nil, verr.Root(err)
		}
		return nil, nil
	}

	return &rtn, err
}

func updateHtml(r *http.Request, key string) error {

	c := appengine.NewContext(r)

	u, err := getUser(r)
	if err != nil {
		return err
	}

	art, err := updateArticle(r, key, time.Now())
	if err != nil {
		return err
	}

	html, err := getHtml(r, key)
	if err != nil {
		return err
	}

	data := &HtmlData{}
	dk := getHtmlDataKey(r, key)

	//get html
	if html == nil {
		// first
		html = &Html{}
		k := getHtmlKey(r, key)

		html.SetKey(k)
		data.SetKey(dk)
	} else {
		err = ds.Get(c, dk, data)
		if err != nil {
			return err
		}
	}

	html.Title = art.Title
	html.SubTitle = art.SubTitle

	html.Author = u.Name
	html.AuthorID = u.Key.StringID()

	err = ds.Put(c, html)
	if err != nil {
		return err
	}

	b, err := createHtml(r, art, u, html)
	if err != nil {
		return err
	}
	data.Content = b

	err = ds.Put(c, data)
	return err
}

func selectHtml(r *http.Request, p int) ([]Html, error) {

	c := appengine.NewContext(r)
	item, err := memcache.Get(c, "html_"+p+"_cursor")
	cursor := ""
	if err == nil {
		cursor = string(item.Value)
	}

	q := datastore.NewQuery(KIND_HTML).
		Order("- CreatedAt").
		Limit(5)

	if cursor != "" {
		cur, err := datastore.DecodeCursor(cursor)
		if err == nil {
			q = q.Start(cur)
		}
	}

	var s []Html

	t := q.Run(c)
	for {
		var h Html
		key, err := t.Next(&h)

		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		h.SetKey(key)
		s = append(s, h)
	}

	cur, err := t.Cursor()
	if err != nil {
		return nil, err
	}

	err = memcache.Set(c, &memcache.Item{
		Key:   "html_" + strconv.Itoa(p+1) + "_cursor",
		Value: []byte(cur.String()),
	})

	if err != nil {
		return nil, err
	}

	return s, nil
}

func deleteHtml(r *http.Request, id string) error {

	c := appengine.NewContext(r)

	hkey := getHtmlKey(r, id)
	err := ds.Delete(c, hkey)
	if err != nil {
		return err
	}
	hdkey := getHtmlDataKey(r, id)
	err = ds.Delete(c, hdkey)
	return err
}

const KIND_HTMLDATA = "HtmlData"

type HtmlData struct {
	key     *datastore.Key
	Content datastore.ByteString `datastore:",noindex"`
}

func getHtmlDataKey(r *http.Request, key string) *datastore.Key {
	c := appengine.NewContext(r)
	return datastore.NewKey(c, KIND_HTMLDATA, key, 0, nil)
}

func (d *HtmlData) GetKey() *datastore.Key {
	return d.key
}

func (d *HtmlData) SetKey(k *datastore.Key) {
	d.key = k
}

func getHtmlData(r *http.Request, k string) (*HtmlData, error) {

	c := appengine.NewContext(r)
	rtn := HtmlData{}
	key := getHtmlDataKey(r, k)
	err := ds.Get(c, key, &rtn)

	if err != nil {
		if verr.Root(err) == datastore.ErrNoSuchEntity {
			return nil, err
		} else if verr.Root(err) == datastore.ErrNoSuchEntity {
			return nil, nil
		}
	}
	return &rtn, nil
}

const KIND_FILE = "File"

type File struct {
	Size int64
	Type int64
	ds.Meta
}

const (
	FILE_TYPE_BG     = 1
	FILE_TYPE_AVATAR = 2
	FILE_TYPE_DATA   = 3
)

func getFileKey(r *http.Request, name string) *datastore.Key {
	c := appengine.NewContext(r)
	return datastore.NewKey(c, KIND_FILE, name, 0, nil)
}

func selectFile(r *http.Request, p int) ([]File, error) {

	c := appengine.NewContext(r)

	item, err := memcache.Get(c, "file_"+strconv.Itoa(p)+"_cursor")
	cursor := ""
	if err == nil {
		cursor = string(item.Value)
	}

	q := datastore.NewQuery(KIND_FILE).
		Filter("Type =", 3).
		Order("- UpdatedAt").
		Limit(10)

	if cursor != "" {
		cur, err := datastore.DecodeCursor(cursor)
		if err == nil {
			q = q.Start(cur)
		}
	}

	var s []File

	t := q.Run(c)
	for {
		var f File
		key, err := t.Next(&f)

		if err == datastore.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		f.SetKey(key)
		s = append(s, f)
	}

	cur, err := t.Cursor()
	if err != nil {
		return nil, err
	}

	err = memcache.Set(c, &memcache.Item{
		Key:   "file_" + strconv.Itoa(p+1) + "_cursor",
		Value: []byte(cur.String()),
	})

	if err != nil {
		return nil, err
	}

	return s, nil
}

func deleteFile(r *http.Request, id string) error {

	c := appengine.NewContext(r)

	fkey := getFileKey(r, id)
	err := ds.Delete(c, fkey)
	if err != nil {
		return err
	}

	fdkey := getFileDataKey(r, id)
	err = ds.Delete(c, fdkey)
	return err
}

func existsFile(r *http.Request, id string, t int64) (bool, error) {

	c := appengine.NewContext(r)
	dir := "data"
	if t == FILE_TYPE_BG {
		dir = "bg"
	} else if t == FILE_TYPE_AVATAR {
		dir = "avatar"
	}

	key := getFileKey(r, dir+"/"+id)

	rtn := File{}
	err := ds.Get(c, key, &rtn)
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return true, verr.Root(err)
		} else {
			return false, nil
		}
	}

	return true, nil

}

func saveFile(r *http.Request, id string, t int64) error {

	upload, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer upload.Close()

	b, err := convertImage(upload)
	if err != nil {
		return err
	}

	c := appengine.NewContext(r)

	dir := "data"
	if t == FILE_TYPE_BG {
		dir = "bg"
	} else if t == FILE_TYPE_AVATAR {
		dir = "avatar"
	} else {
		if id == "" {
			id = header.Filename
		}
	}

	fid := dir + "/" + id
	file := &File{
		Size: int64(len(b)),
		Type: t,
	}

	file.Key = getFileKey(r, fid)
	err = ds.Put(c, file)
	if err != nil {
		return err
	}
	fileData := &FileData{
		Content: b,
		Mime:    header.Header["Content-Type"][0],
	}
	fileData.SetKey(getFileDataKey(r, fid))
	err = ds.Put(c, fileData)
	if err != nil {
		return err
	}
	return nil
}

func saveAvatar(r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)
	err := saveFile(r, u.ID, FILE_TYPE_AVATAR)
	return err
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
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return nil, verr.Root(err)
		} else if verr.Root(err) == datastore.ErrNoSuchEntity {
			return nil, nil
		}
	}

	return &rtn, nil
}
