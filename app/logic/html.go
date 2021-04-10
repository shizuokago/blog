package logic

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/shizuokago/blog/datastore"

	"golang.org/x/tools/present"
	"golang.org/x/xerrors"
)

//go:embed _entry
var embEntry embed.FS
var tmpl *template.Template

func init() {
	var err error
	tmpl, err = createTemplate()
	if err != nil {
		log.Println(err)
	}
}

func createTemplate() (*template.Template, error) {

	var err error

	action := "_entry/action.tmpl"
	entry := "_entry/entry.tmpl"

	funcMap := template.FuncMap{
		"playable": playable,
		"convert":  convert,
	}

	tmpl = present.Template()
	tmpl = tmpl.Funcs(funcMap)

	tmpl, err = tmpl.ParseFS(embEntry, action, entry)
	if err != nil {
		return nil, xerrors.Errorf("action template parse: %w", err)
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

func CreateHTML(ctx context.Context, id string, blog *datastore.Blog, art *datastore.Article, u *datastore.User) (*datastore.HTMLParam, error) {

	//create header
	header := art.Title + "\n\n" +
		u.Name + "\n" +
		u.Job + "\n" +
		u.Email + "\n" +
		u.URL + "\n" +
		u.TwitterId + "\n"

	mark := string(art.Markdown)
	txt := header + "\n" + mark

	desc := strings.ReplaceAll(mark, "\n", "")
	if len(desc) > 600 {
		desc = string([]rune(desc)[:200]) + "..."
	}

	ds := FileDs{
		ctx: ctx,
	}

	pctx := present.Context{ReadFile: ds.readFile}

	reader := strings.NewReader(txt)
	doc, err := pctx.Parse(reader, "blog.article", 0)
	if err != nil {
		return nil, xerrors.Errorf("context parse: %w", err)
	}

	html := &datastore.HTML{}

	html.Title = art.Title
	html.SubTitle = art.SubTitle
	// first
	html.Author = u.Name
	html.AuthorID = u.Email
	html.Updater = u.Name
	html.UpdaterID = u.Email

	rtn := struct {
		*present.Doc
		Template    *template.Template
		PlayEnabled bool
		StringID    string
		Description string
		Blog        *datastore.Blog
		HTML        *datastore.HTML
	}{doc, tmpl, true, id, desc, blog, html}

	//Render
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err = tmpl.ExecuteTemplate(writer, "root", rtn)

	if err != nil {
		return nil, xerrors.Errorf("execute template: %w", err)
	}
	writer.Flush()

	data := &datastore.HTMLData{}
	data.Content = b.Bytes()

	p := datastore.HTMLParam{}

	p.Body = html
	p.Data = data

	return &p, nil
}
