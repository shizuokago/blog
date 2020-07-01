package main

//Editor Generate
//

import (
	"bufio"
	"bytes"
	"fmt"
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

	present.Register("picture", parsePicture)
}

func main() {

	jQuery(DOCUMENT).Ready(func() {
		draw()
		resize()
	})

	jQuery(js.Global).Resize(func(e jquery.Event) {
		resize()
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

	jQuery("#saveBGBtn").On(jquery.CLICK, func(e jquery.Event) {
		jQuery("#file").Call("click")
	})

	jQuery("#file").On(jquery.CHANGE, func(e jquery.Event) {
		jQuery("#bgForm").Call("submit")
	})

	auto := jQuery("#AutoSave").Val()

	rd, aj := 0, 0
	go func() {
		t := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-t.C:
				rd = rd + 1
				aj = aj + 1
				flag := false

				if rd == 2 {
					flag = redraw()
					rd = 0
				}

				if auto != "" {
					if aj == 4 {
						if flag {
							ajax("autosave")
						}
						aj = 0
					}
				}
			}
		}
		t.Stop()
	}()
}

func ajax(url string) {

	u := url
	if url == "autosave" {
		u = "save"
	}

	var d *js.Object
	if url != "autosave" {
		d = js.Global.Call("waitDialog")
	}

	id := jQuery(ARTICLE_ID).Val()
	data := js.M{
		"Title":    jQuery(TITLE).Val(),
		"Tags":     jQuery(TAGS).Val(),
		"Markdown": jQuery(INPUT).Val(),
	}

	ajaxopt := js.M{
		"async":    true,
		"type":     "POST",
		"url":      "/admin/article/" + u + "/" + id,
		"dataType": "json",
		"data":     data,
		"success": func(data map[string]interface{}) {
		},
		"error": func(status interface{}) {
		},
		"complete": func(status interface{}) {
			if url != "autosave" {
				d.Call("close")
			} else {
				d = js.Global.Call("toast", "Auto Saved")
			}
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
	AuthorID  string
	Updater   string
	UpdaterID string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func draw() {
	h := getHtml()

	doc := js.Global.Get("document")
	iframe := doc.Call("getElementById", "result")

	if iframe != nil {
		idoc := iframe.Get("contentDocument")

		idoc.Call("open")
		idoc.Call("write", h)
		idoc.Call("close")

		unbind()
	}
}

var beforeBody = ""

func redraw() bool {

	h := getHtml()
	r := strings.NewReader(h)

	doc, _ := goquery.NewDocumentFromReader(r)
	body := doc.Find("body")

	bh, _ := body.Html()

	if beforeBody != bh {
		jQuery(OUTPUT).Contents().Find("body").SetHtml(bh)
		unbind()
		beforeBody = bh
		return true
	}
	return false
}

func unbind() {
	iframe := jQuery("#result").Contents()
	iframe.Find("a").RemoveAttr("href")
	iframe.Find("button").RemoveAttr("onclick")
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

	zero := time.Time{}
	h := Html{
		Author:    author,
		AuthorID:  "empty",
		Updater:   author,
		UpdaterID: "empty",
		CreatedAt: zero,
		UpdatedAt: zero,
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
		StringID    string
		Blog        interface{}
		HTML        Html
	}{doc, gblTmpl, true, jQuery(ARTICLE_ID).Val(), blog, h}

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

	file := "/file/data/" + name
	code := "error"
	ajaxopt := js.M{
		"async":    false,
		"type":     "GET",
		"url":      file,
		"dataType": "text/plain",
		"complete": func(status map[string]interface{}) {
			code = status["responseText"].(string)
		},
	}
	jquery.Ajax(ajaxopt)
	return []byte(code), nil
}

func parsePicture(ctx *present.Context, fileName string, lineno int, text string) (present.Elem, error) {

	args := strings.Fields(text)
	img := present.Image{URL: "/file/data/" + args[1]}
	a, err := parseArgs(fileName, lineno, args[2:])
	if err != nil {
		return nil, err
	}
	switch len(a) {
	case 0:
	case 2:
		if v, ok := a[0].(int); ok {
			img.Height = v
		}
		if v, ok := a[1].(int); ok {
			img.Width = v
		}
	default:
		return nil, fmt.Errorf("incorrect image invocation: %q", text)
	}
	return img, nil

}

func parseArgs(name string, line int, args []string) (res []interface{}, err error) {
	res = make([]interface{}, len(args))
	for i, v := range args {
		if len(v) == 0 {
			return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
		}
		switch v[0] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
			}
			res[i] = n
		case '/':
			if len(v) < 2 || v[len(v)-1] != '/' {
				return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
			}
			res[i] = v
		case '$':
			res[i] = "$"
		case '_':
			if len(v) == 1 {
				// Do nothing; "_" indicates an intentionally empty parameter.
				break
			}
			fallthrough
		default:
			return nil, fmt.Errorf("%s:%d bad code argument %q", name, line, v)
		}
	}
	return
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
