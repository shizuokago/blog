package main

import (
	"context"
	"fmt"
	"strings"
	"syscall/js"
	"time"

	"github.com/shizuokago/blog/datastore"
	"github.com/shizuokago/blog/logic"
)

var saveMD string

func main() {

	done := make(chan struct{}, 0)
	fmt.Println("wasm initialize")
	err := run()
	if err != nil {
		fmt.Printf("wasm error: %+v", err)
	}
	fmt.Println("wasm set success")
	<-done
}

func run() error {

	/*
		doc := js.Global().Get("document")
		win := js.Global().Get("window")

		doc.Call(
			"addEventListener",
			"DOMContentLoaded",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				fmt.Println("ContentLoaded() js")
				return nil
			}),
		)
		win.Call(
			"addEventListener",
			"resize",
			js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				draw()
				resize()
				return nil
			}),
		)
	*/

	draw()
	for range time.Tick(5 * time.Second) {
		draw()
	}
	return nil
}

func getID(id string) js.Value {
	return js.Global().Get(Document).Call("getElementById", id)
}

func draw() {

	md := getID(INPUT).Get("value").String()

	if md == saveMD {
		return
	}
	saveMD = md

	h := getHTML(md)

	doc := js.Global().Get("document")
	iframe := doc.Call("getElementById", "result")

	if !iframe.IsUndefined() {
		idoc := iframe.Get("contentDocument")
		idoc.Call("open")
		idoc.Call("write", h)
		idoc.Call("close")
	}
}

func getHTML(md string) string {

	path := js.Global().Get("location").Get("pathname")

	paths := strings.Split(path.String(), "/")
	id := paths[len(paths)-1]

	title := getID(TITLE).Get("value").String()
	tags := getID(TAGS).Get("value").String()

	author := getID(AUTHOR).Get("value").String()
	job := getID(JOB).Get("value").String()
	mail := getID(EMAIL).Get("value").String()
	url := getID(URL).Get("value").String()
	twitter := "@" + getID(TWITTER).Get("value").String()

	art := datastore.Article{}
	user := datastore.User{}

	art.Title = title
	art.SubTitle = title
	art.Tags = tags
	art.PublishDate = time.Now()

	art.Markdown = []byte(md)

	user.Name = author
	user.Job = job
	user.Email = mail
	user.URL = url
	user.TwitterId = twitter
	//user.Meta

	blog := datastore.Blog{}

	p, err := logic.CreateHTML(context.Background(), id, &blog, &art, &user)
	if err != nil {
		return err.Error()
	}

	return string(p.Data.Content)
}

const (
	TITLE      = "Title"
	TAGS       = "Tags"
	AUTHOR     = "Name"
	JOB        = "Job"
	EMAIL      = "Email"
	URL        = "URL"
	TWITTER    = "TwitterId"
	ARTICLE_ID = "ID"
	BLOGNAME   = "BlogName"
	Document   = "document"
	SAVE       = "save"
	PUBLISH    = "publish"
	OUTPUT     = "result"

	INPUT = "editor"
	LEFT  = "left"
	RIGHT = "right"
)
