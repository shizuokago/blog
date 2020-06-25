package datastore

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/datastore"
	"golang.org/x/tools/present"
	"google.golang.org/api/iterator"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = createTemplate()
	if err != nil {
		panic(err)
	}
}

func CreateHtml(r *http.Request, art *Article, u *User, html *Html) ([]byte, error) {

	//create header
	header := art.Title + "\n\n" +
		u.Name + "\n" +
		u.Job + "\n" +
		u.Email + "\n" +
		u.URL + "\n" +
		u.TwitterId + "\n"

	txt := header + "\n" + string(art.Markdown)

	ds := FileDs{
		request: r,
	}
	ctx := present.Context{ReadFile: ds.readFile}

	reader := strings.NewReader(txt)
	doc, err := ctx.Parse(reader, "blog.article", 0)
	if err != nil {
		return nil, err
	}

	bgd := GetBlog(r)
	rtn := struct {
		*present.Doc
		Template    *template.Template
		PlayEnabled bool
		StringID    string
		Blog        *Blog
		HTML        *Html
	}{doc, tmpl, true, art.Key.Name, bgd, html}

	//Render
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = tmpl.ExecuteTemplate(writer, "root", rtn)

	if err != nil {
		return nil, err
	}
	writer.Flush()

	return b.Bytes(), nil
}

func createTemplate() (*template.Template, error) {

	action := "./cmd/static/templates/entry/action.tmpl"
	entry := "./cmd/static/templates/entry/entry.tmpl"

	tmpl = present.Template()
	funcMap := template.FuncMap{
		"playable": playable,
		"convert":  convert,
	}
	tmpl = tmpl.Funcs(funcMap)
	_, err := tmpl.ParseFiles(action, entry)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}

func convert(t time.Time) string {
	if t.IsZero() {
		return "None"
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	jt := t.In(jst)
	return jt.Format("2006/01/02 15:04")
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

	c := r.Context()

	rtn := Html{}
	key := getHtmlKey(r, k)

	client, err := createClient(c)
	err = client.Get(c, key, &rtn)

	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		}
		return nil, err
	}

	return &rtn, err
}

func UpdateHtml(r *http.Request, mail, key string) error {

	c := r.Context()

	u, err := GetUser(r, mail)
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
	data.Content = b

	_, err = client.Put(c, data.GetKey(), data)
	return err
}

func SelectHtml(r *http.Request, p int) ([]Html, error) {

	c := r.Context()

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

	c := r.Context()

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
	Content []byte `datastore:",noindex"`
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

	c := r.Context()
	rtn := HtmlData{}
	key := getHtmlDataKey(r, k)

	client, err := createClient(c)
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
