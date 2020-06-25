package datastore

import (
	"context"
	"net/http"
	"strings"
	"time"

	verr "github.com/knightso/base/errors"
	"github.com/pborman/uuid"

	"cloud.google.com/go/datastore"

	"google.golang.org/api/iterator"
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

func init() {
}

func createClient(ctx context.Context) (*datastore.Client, error) {
	client, err := datastore.NewClient(ctx, "shizuoka-go")
	return client, err
}

const KIND_BLOG = "Blog"

type Blog struct {
	Name        string
	Author      string
	Tags        string
	Description string
	Template    string
	Meta
}

var pkgBlog = Blog{}

func GetBlog(r *http.Request) *Blog {

	if pkgBlog.Name != "" {
		return &pkgBlog
	}
	c := appengine.NewContext(r)
	key := datastore.NameKey(KIND_BLOG, "Fixing", nil)

	client, err := createClient(c)
	if err != nil {
		//
	}

	err = client.Get(c, key, &pkgBlog)
	if err != nil {
		// Nothing
	}
	return &pkgBlog
}

func PutBlog(r *http.Request) error {

	pkgBlog = Blog{
		Name:        r.FormValue("BlogName"),
		Author:      r.FormValue("BlogAuthor"),
		Description: r.FormValue("Description"),
		Tags:        r.FormValue("BlogTags"),
		Template:    r.FormValue("BlogTemplate"),
	}

	c := appengine.NewContext(r)
	key := datastore.NameKey(KIND_BLOG, "Fixing", nil)

	pkgBlog.SetKey(key)

	client, err := createClient(c)

	_, err = client.Put(c, key, &pkgBlog)
	if err != nil {
		return err
	}

	return nil
}

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

func getUserKey(r *http.Request) *datastore.Key {
	c := appengine.NewContext(r)
	u := user.Current(c)

	return datastore.NameKey(KIND_USER, u.ID, nil)
}

func GetUser(r *http.Request) (*User, error) {

	c := appengine.NewContext(r)

	rtn := User{}
	key := getUserKey(r)

	client, err := createClient(c)

	err = client.Get(c, key, &rtn)
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return nil, verr.Root(err)
		} else {
			return nil, nil
		}
	}
	return &rtn, nil
}

func PutInformation(r *http.Request) (*User, error) {

	c := appengine.NewContext(r)

	r.ParseForm()

	save := false
	if r.FormValue("AutoSave") != "" {
		save = true
	}

	rtn := User{
		Name:      r.FormValue("Name"),
		Job:       r.FormValue("Job"),
		Email:     r.FormValue("Email"),
		URL:       r.FormValue("Url"),
		TwitterId: r.FormValue("TwitterId"),
		AutoSave:  save,
	}

	err := PutBlog(r)
	if err != nil {
		return nil, err
	}

	//function
	rtn.Key = getUserKey(r)

	client, err := createClient(c)

	_, err = client.Put(c, rtn.Key, &rtn)
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
	Markdown    string `datastore:",noindex"`
	Meta
}

func getArticleKey(r *http.Request, id string) *datastore.Key {
	return datastore.NameKey(KIND_ARTICLE, id, nil)
}

func SelectArticle(r *http.Request, p int) ([]Article, error) {

	c := appengine.NewContext(r)

	//TODO CURCOR

	q := datastore.NewQuery("Article").Order("- UpdatedAt").Limit(10)

	var s []Article

	client, err := createClient(c)
	if err != nil {
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
	c := appengine.NewContext(r)

	rtn := Article{}
	key := getArticleKey(r, id)

	client, err := createClient(c)

	err = client.Get(c, key, &rtn)
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return nil, verr.Root(err)
		} else if verr.Root(err) == datastore.ErrNoSuchEntity {
			return nil, nil
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

	c := appengine.NewContext(r)

	art.Title = title
	art.SubTitle = CreateSubTitle(r.FormValue("Markdown"))
	art.Tags = tags
	art.Markdown = mark
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

func CreateSubTitle(src string) string {

	dst := strings.Replace(src, "\n", "", -1)
	dst = strings.Replace(dst, "*", "", -1)

	if len(dst) > 600 {
		dst = string([]rune(dst)[0:200]) + "..."
	}
	return dst
}

func CreateArticle(r *http.Request) (string, error) {

	c := appengine.NewContext(r)
	id := uuid.New()

	bgd := GetBlog(r)
	base := bgd.Template
	article := &Article{
		Title:    "New Title",
		Tags:     bgd.Tags,
		Markdown: base,
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

func CreateHtmlFromMail(r *http.Request, d *MailData) error {

	c := appengine.NewContext(r)
	id := uuid.New()

	bgd := GetBlog(r)
	article := &Article{
		Title:    d.subject,
		Tags:     bgd.Tags,
		Markdown: string(d.msg.Bytes()),
	}

	article.Key = getArticleKey(r, id)
	client, err := createClient(c)
	_, err = client.Put(c, article.Key, article)
	if err != nil {
		return err
	}

	fid := "bg/" + id
	fb := d.file.Bytes()
	file := &File{
		Size: int64(len(fb)),
		Type: FILE_TYPE_BG,
	}

	file.Key = getFileKey(r, fid)

	_, err = client.Put(c, file.Key, file)
	if err != nil {
		return err
	}
	fileData := &FileData{
		Content: fb,
		Mime:    d.mime,
	}

	fdk := getFileDataKey(r, fid)
	fileData.SetKey(fdk)
	_, err = client.Put(c, fileData.GetKey(), fileData)
	if err != nil {
		return err
	}
	return nil

}

func DeleteArticle(r *http.Request, id string) error {

	c := appengine.NewContext(r)

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

const KIND_HTML = "Html"

type Html struct {
	Title     string
	SubTitle  string
	Author    string
	AuthorID  string
	Updater   string
	UpdaterID string
	Meta
}

func getHtmlKey(r *http.Request, key string) *datastore.Key {
	return datastore.NameKey(KIND_HTML, key, nil)
}

func GetHtml(r *http.Request, k string) (*Html, error) {

	c := appengine.NewContext(r)

	rtn := Html{}
	key := getHtmlKey(r, k)

	client, err := createClient(c)
	err = client.Get(c, key, &rtn)

	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return nil, verr.Root(err)
		}
		return nil, nil
	}

	return &rtn, err
}

func UpdateHtml(r *http.Request, key string) error {

	c := appengine.NewContext(r)

	u, err := GetUser(r)
	if err != nil {
		return err
	}

	art, err := UpdateArticle(r, key, time.Now())
	if err != nil {
		return err
	}

	html, err := GetHtml(r, key)
	if err != nil {
		return err
	}

	data := &HtmlData{}
	dk := getHtmlDataKey(r, key)

	client, err := createClient(c)

	//get html
	if html == nil {
		// first
		html = &Html{}
		k := getHtmlKey(r, key)

		html.SetKey(k)
		data.SetKey(dk)

		html.Author = u.Name
		html.AuthorID = u.Key.Name
	} else {
		err = client.Get(c, dk, data)
		if err != nil {
			return err
		}
		html.Updater = u.Name
		html.UpdaterID = u.Key.Name
	}

	html.Title = art.Title
	html.SubTitle = art.SubTitle

	_, err = client.Put(c, html.GetKey(), html)
	if err != nil {
		return err
	}

	b, err := CreateHtml(r, art, u, html)
	if err != nil {
		return err
	}
	data.Content = string(b)

	_, err = client.Put(c, data.GetKey(), data)
	return err
}

func SelectHtml(r *http.Request, p int) ([]Html, error) {

	c := appengine.NewContext(r)

	q := datastore.NewQuery(KIND_HTML).
		Order("- CreatedAt").
		Limit(5)

	var s []Html

	client, err := createClient(c)
	if err != nil {
		return nil, err
	}

	t := client.Run(c, q)
	for {
		var h Html
		key, err := t.Next(&h)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		h.SetKey(key)
		s = append(s, h)
	}

	return s, nil
}

func DeleteHtml(r *http.Request, id string) error {

	c := appengine.NewContext(r)

	hkey := getHtmlKey(r, id)

	client, err := createClient(c)

	err = client.Delete(c, hkey)
	if err != nil {
		return err
	}
	hdkey := getHtmlDataKey(r, id)
	err = client.Delete(c, hdkey)
	return err
}

const KIND_HTMLDATA = "HtmlData"

type HtmlData struct {
	key     *datastore.Key
	Content string `datastore:",noindex"`
}

func getHtmlDataKey(r *http.Request, key string) *datastore.Key {
	return datastore.NameKey(KIND_HTMLDATA, key, nil)
}

func (d *HtmlData) GetKey() *datastore.Key {
	return d.key
}

func (d *HtmlData) SetKey(k *datastore.Key) {
	d.key = k
}

func GetHtmlData(r *http.Request, k string) (*HtmlData, error) {

	c := appengine.NewContext(r)
	rtn := HtmlData{}
	key := getHtmlDataKey(r, k)

	client, err := createClient(c)
	err = client.Get(c, key, &rtn)

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
	Meta
}

const (
	FILE_TYPE_BG     = 1
	FILE_TYPE_AVATAR = 2
	FILE_TYPE_DATA   = 3
)

func getFileKey(r *http.Request, name string) *datastore.Key {
	return datastore.NameKey(KIND_FILE, name, nil)
}

func SelectFile(r *http.Request, p int) ([]File, error) {

	c := appengine.NewContext(r)

	q := datastore.NewQuery(KIND_FILE).
		Filter("Type =", 3).
		Order("- UpdatedAt").
		Limit(10)

	var s []File

	client, err := createClient(c)
	if err != nil {
		return nil, err
	}
	t := client.Run(c, q)
	for {
		var f File
		key, err := t.Next(&f)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		f.SetKey(key)
		s = append(s, f)
	}

	return s, nil
}

func DeleteFile(r *http.Request, id string) error {

	c := appengine.NewContext(r)

	fkey := getFileKey(r, id)

	client, err := createClient(c)
	err = client.Delete(c, fkey)
	if err != nil {
		return err
	}

	fdkey := getFileDataKey(r, id)
	err = client.Delete(c, fdkey)
	return err
}

func ExistsFile(r *http.Request, id string, t int64) (bool, error) {

	c := appengine.NewContext(r)
	dir := "data"
	if t == FILE_TYPE_BG {
		dir = "bg"
	} else if t == FILE_TYPE_AVATAR {
		dir = "avatar"
	}

	key := getFileKey(r, dir+"/"+id)

	rtn := File{}
	client, err := createClient(c)
	err = client.Get(c, key, &rtn)
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return true, verr.Root(err)
		} else {
			return false, nil
		}
	}

	return true, nil

}

func SaveFile(r *http.Request, id string, t int64) error {

	upload, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer upload.Close()

	b, flg, err := ConvertImage(upload)
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

	client, err := createClient(c)
	_, err = client.Put(c, file.Key, file)
	if err != nil {
		return err
	}

	mime := header.Header["Content-Type"][0]
	if flg {
		mime = "image/jpeg"
	}

	fileData := &FileData{
		Content: b,
		Mime:    mime,
	}
	fileData.SetKey(getFileDataKey(r, fid))
	_, err = client.Put(c, fileData.GetKey(), fileData)
	if err != nil {
		return err
	}
	return nil
}

func SaveAvatar(r *http.Request) error {
	c := appengine.NewContext(r)
	u := user.Current(c)
	err := SaveFile(r, u.ID, FILE_TYPE_AVATAR)
	return err
}

func SaveBackgroundImage(r *http.Request, id string) error {
	err := SaveFile(r, id, FILE_TYPE_BG)
	return err
}

func DeleteBackgroundImage(r *http.Request, id string) error {
	err := DeleteFile(r, "bg/"+id)
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
	return datastore.NameKey(KIND_FILEDATA, name, nil)
}

func GetFileData(r *http.Request, name string) (*FileData, error) {
	c := appengine.NewContext(r)

	rtn := FileData{}
	key := getFileDataKey(r, name)

	client, err := createClient(c)
	err = client.Get(c, key, &rtn)
	if err != nil {
		if verr.Root(err) != datastore.ErrNoSuchEntity {
			return nil, verr.Root(err)
		} else if verr.Root(err) == datastore.ErrNoSuchEntity {
			return nil, nil
		}
	}

	return &rtn, nil
}
