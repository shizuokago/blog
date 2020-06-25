package datastore

import (
	"bufio"
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"time"

	"golang.org/x/tools/present"
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
	}{doc, tmpl, true, art.Key.StringID(), bgd, html}

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

	action := "templates/entry/action.tmpl"
	entry := "templates/entry/entry.tmpl"

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
