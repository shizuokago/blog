package main

//Editor Generate
//

import (
	"bufio"
	"bytes"
	"html/template"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/jquery"

	"golang.org/x/tools/present"
)

var gblTmpl *template.Template
var jQuery = jquery.NewJQuery

func init() {
	gblTmpl = present.Template()
	funcMap := template.FuncMap{
		"playable": playable,
		"convert":  convert,
	}
	_, err := gblTmpl.Funcs(funcMap).Parse(TMPL)
	if err != nil {
		panic(err)
	}
}

func main() {

	jQuery(DOCUMENT).Ready(func() {
		draw()
		resize()
	})

	jQuery(js.Global).Resize(func(e jquery.Event) {
		resize()
	})

	cnt := 0
	jQuery(INPUT).On(jquery.KEYDOWN, func(e jquery.Event) {
		cnt++
		if cnt == 15 {
			redraw()
			cnt = 0
		}
	})

	jQuery(PUBLISH).On(jquery.CLICK, func(e jquery.Event) {
		ajax("publish")
	})

	jQuery(SAVE).On(jquery.CLICK, func(e jquery.Event) {
		ajax("save")
	})

	jQuery("button#delete").On(jquery.CLICK, func(e jquery.Event) {
		url := "/admin/article/delete/" + jQuery(ARTICLE_ID).Val()
		l := js.Global.Get("location")
		l.Set("href", url)
	})

	jQuery("button#private").On(jquery.CLICK, func(e jquery.Event) {
		ajax("private")
	})

	jQuery("button#viewBtn").On(jquery.CLICK, func(e jquery.Event) {
		url := "/entry/" + jQuery(ARTICLE_ID).Val()
		js.Global.Call("open", url, "_blank")
	})

	jQuery("#file").On(jquery.CLICK, func(e jquery.Event) {
		jQuery("#file").Call("click")
	})

	jQuery("#file").On(jquery.CHANGE, func(e jquery.Event) {
		jQuery("#bgForm").Call("submit")
	})
}

func ajax(url string) {

	d := js.Global.Get("document")
	dialog := d.Call("querySelector", "#wait_dialog")
	dialog.Call("showModal")

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
		"complete": func(status interface{}) {
			dialog.Call("close")
		},
	}

	jquery.Ajax(ajaxopt)
}

func resize() {

	height := jQuery(js.Global).Height()

	margin := 230

	out := jQuery(OUTPUT)
	empty := jQuery("img#empty")

	jQuery(LEFT).SetHeight(strconv.Itoa(height - margin))
	jQuery(RIGHT).SetHeight(strconv.Itoa(height - margin))

	jQuery(INPUT).SetHeight(strconv.Itoa(height - margin))
	out.SetHeight(strconv.Itoa(height - margin))

	empty.SetCss("top", out.Css("top"))
	empty.SetCss("left", out.Css("left"))
	empty.SetCss("width", out.Width()-25)
	empty.SetCss("height", out.Height())
}

type Html struct {
	Author    string
	CreatedAt time.Time
}

func draw() {
	h := getHtml()

	doc := js.Global.Get("document")
	iframe := doc.Call("getElementById", "result")
	idoc := iframe.Get("contentDocument")

	idoc.Call("open")
	idoc.Call("write", h)
	idoc.Call("close")
}

func redraw() {

	h := getHtml()
	r := strings.NewReader(h)

	doc, _ := goquery.NewDocumentFromReader(r)
	body := doc.Find("body")

	bh, _ := body.Html()
	jQuery(OUTPUT).Contents().Find("body").SetHtml(bh)
}

func getHtml() string {

	title := jQuery(TITLE).Val()

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
		return err.Error()
	}

	h := Html{
		Author:    author,
		CreatedAt: time.Now(),
	}

	blog := struct {
		Name        string
		Author      string
		Tags        string
		Description string
	}{"Empty", "Empty", "Empty", "Empty"}

	rtn := struct {
		*present.Doc
		Template    *template.Template
		PlayEnabled bool
		AuthorID    string
		StringID    string
		Blog        interface{}
		HTML        Html
	}{doc, gblTmpl, true, "empty", jQuery(ARTICLE_ID).Val(), blog, h}

	//Render
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	err = gblTmpl.ExecuteTemplate(writer, "root", rtn)
	if err != nil {
		return err.Error()
	}
	writer.Flush()

	return string(b.Bytes())
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}

func convert(t time.Time) string {
	//jst, _ := time.LoadLocation("Asia/Tokyo")
	//jt := t.In(jst)
	return t.Format("2006/01/02 15:04")
}

func readFile(name string) ([]byte, error) {

	loc := js.Global.Get("location")
	//get host
	host := loc.Get("origin")
	file := "/file/data/" + name

	url := host.String() + file

	// request
	//resp, err := http.Get(url)
	//if err != nil {
	//return nil, err
	//}
	//defer resp.Body.Close()

	//byteArray, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//return nil, err
	//}

	return []byte(url), nil
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
	BLOGNAME   = "input#BlogName"
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
