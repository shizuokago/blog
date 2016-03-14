package main

//Editor Generate
//

import (
	"bufio"
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

	jQuery(PUBLISH).On(jquery.CLICK, func(e jquery.Event) {
		ajax("publish")
	})

	jQuery(SAVE).On(jquery.CLICK, func(e jquery.Event) {
		ajax("save")
	})
}

func ajax(url string) {
	id := jQuery(ARTICLE_ID).Val()
	data := js.M{
		"Title":    jQuery(TITLE).Val(),
		"Tags":     jQuery(TAGS).Val(),
		"Markdown": jQuery(INPUT).Val(),
	}
	ajaxopt := js.M{
		"async":    true,
		"type":     "POST",
		"url":      "/admin/article/" + url + "/" + id,
		"dataType": "json",
		"data":     data,
		"success": func(data map[string]interface{}) {
		},
		"error": func(status interface{}) {
		},
	}
	//ajax call:
	jquery.Ajax(ajaxopt)
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
	ctx := present.Context{ReadFile: readFile}

	reader := strings.NewReader(art)
	doc, err := ctx.Parse(reader, "blog.article", 0)
	if err != nil {
		return
	}

	rtn := struct {
		*present.Doc
		Template    *template.Template
		PlayEnabled bool
		StringID    string
	}{doc, gblTmpl, true, jQuery(ARTICLE_ID).Val()}
	//Render
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = gblTmpl.ExecuteTemplate(writer, "root", rtn)
	if err != nil {
		return
	}
	writer.Flush()

	jQuery(OUTPUT).Contents().Find("html").SetHtml(string(b.Bytes()))
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}

func readFile(name string) ([]byte, error) {
	//select file data
	return nil, nil
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
	SAVE       = "button#save"
	PUBLISH    = "button#publish"
	OUTPUT     = "iframe#result"
	LEFT       = "div#left"
	RIGHT      = "div#right"
	TMPL       = `
	{{ .ARTICLE_TMPL }}
	{{ .ACTION_TMPL }}
`
)
