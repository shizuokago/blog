package main

//Editor Generate
//
//  # loading template(article.tmpl,action.tmpl)
//
//  # gopherjs build editor.go
//  # mv editor*js* ../static/js
//

import (
	"bytes"
	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"
	"golang.org/x/tools/present"
	"html/template"
	"strconv"
	"strings"
)

var gblTmpl *template.Template
var jQuery = jquery.NewJQuery

func init() {
	gblTmpl = present.Template()
	gblTmpl.Funcs(template.FuncMap{"playable": playable})
	_, err := gblTmpl.Parse(TMPL)
	if err != nil {
		panic(err)
	}
}

func parseArticle(article string) (*present.Doc, error) {
	r := strings.NewReader(article)
	return present.Parse(r, "root", 0)
}

func render(doc *present.Doc) (*bytes.Buffer, error) {
	w := bytes.NewBuffer([]byte{})
	err := doc.Render(w, gblTmpl)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func main() {

	jQuery(DOCUMENT).Ready(func() {
		redraw()
		resize()
	})

	jQuery(js.Global).Resize(func(e jquery.Event) {
		resize()
	})

	jQuery(INPUT).On(jquery.KEYDOWN, func(e jquery.Event) {
		redraw()
	})

	jQuery(BUTTON).On(jquery.CLICK, func(e jquery.Event) {

		id := jQuery(ARTICLE_ID).Val()
		//Title,Tags,Markdown
		ajaxopt := map[string]interface{}{
			"async":       true,
			"type":        "POST",
			"url":         "/admin/article/save/" + id,
			"contentType": "application/json charset=utf-8",
			"dataType":    "json",
			"data":        nil,
			"success": func(data map[string]interface{}) {
			},
			"error": func(status interface{}) {
			},
		}
		//ajax call:
		jquery.Ajax(ajaxopt)

	})
}

func resize() {

	height := jQuery(js.Global).Height()

	margin := 215

	jQuery(LEFT).SetHeight(strconv.Itoa(height - margin))
	jQuery(RIGHT).SetHeight(strconv.Itoa(height - margin))

	jQuery(INPUT).SetHeight(strconv.Itoa(height - margin))
	jQuery(OUTPUT).SetHeight(strconv.Itoa(height - margin))
}

func redraw() {

	title := jQuery(TITLE).Val()

	//sub
	//date
	//tags

	author := jQuery(AUTHOR).Val()
	job := jQuery(JOB).Val()
	mail := jQuery(EMAIL).Val()
	url := jQuery(URL).Val()
	twitter := "@" + jQuery(TWITTER).Val()

	header := title + "\n\n" +
		author + "\n" +
		job + "\n" +
		mail + "\n" +
		url + "\n" +
		twitter + "\n"

	md := jQuery(INPUT).Val()

	art := header + "\n" + md
	doc, _ := parseArticle(art)

	w, _ := render(doc)

	jQuery(OUTPUT).Contents().Find("html").SetHtml(w.String())
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}

const (
	TITLE      = "input#Title"
	TAGS       = "input#Tags"
	AUTHOR     = "input#Name"
	JOB        = "input#Job"
	EMAIL      = "input#Email"
	URL        = "input#URL"
	TWITTER    = "input#TwitterId"
	ARTICLE_ID = "input#ID"
	DOCUMENT   = "document"
	INPUT      = "textarea#editor"
	BUTTON     = "button#save"
	OUTPUT     = "iframe#result"
	LEFT       = "div#left"
	RIGHT      = "div#right"
	TMPL       = `
	{{ .ARTICLE_TMPL }}
	{{ .ACTION_TMPL }}
`
)
