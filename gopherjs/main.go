package main

import (
	"bytes"
	"github.com/gopherjs/jquery"
	"golang.org/x/tools/present"
	"html/template"
	"strings"
)

var gblTmpl *template.Template
var jQuery = jquery.NewJQuery

func init() {
	gblTmpl = present.Template()
	gblTmpl.Funcs(template.FuncMap{"playable": playable})
	_, err := gblTmpl.Parse(ARTICLE_TMPL)
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

	jQuery("document").Ready(func() {
		jQuery(INPUT).SetText(ARTICLE)
		redraw()
	})

	jQuery(BUTTON).On(jquery.CLICK, func(e jquery.Event) {
		redraw()
	})
}

func redraw() {

	doc, _ := parseArticle(jQuery(INPUT).Val())

	w, _ := render(doc)
	jQuery(OUTPUT).Contents().Find("html").SetHtml(w.String())
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}

const (
	INPUT        = "textarea#editor"
	BUTTON       = "button#save"
	OUTPUT       = "iframe#result"
	ARTICLE_TMPL = `
{/* This is the article template. It defines how articles are formatted. */}
{{define "root"}}
<!DOCTYPE html>
<html>
  <head>
    <title>{{.Title}}</title>
    <link type="text/css" rel="stylesheet" href="./article.css">
    <meta charset='utf-8'>
  </head>

  <body>
    <div id="topbar" class="wide">
      <div class="container">
        <div id="heading">{{.Title}}
          {{with .Subtitle}}{{.}}{{end}}
        </div>
      </div>
    </div>
    <div id="page" class="wide">
      <div class="container">
        {{with .Sections}}
          <div id="toc">
            {{template "TOC" .}}
          </div>
        {{end}}

        {{range .Sections}}
          {{elem $.Template .}}
        {{end}}{{/* of Section block */}}

        {{if .Authors}}
          <h2>Authors</h2>
          {{range .Authors}}
            <div class="author">
              {{range .Elem}}{{elem $.Template .}}{{end}}
            </div>
          {{end}}
        {{end}}
      </div>
    </div>
<!--
    <script src='/play.js'></script>
-->
  </body>
</html>
{{end}}

{{define "TOC"}}
  <ul>
  {{range .}}
    <li><a href="#TOC_{{.FormattedNumber}}">{{.Title}}</a></li>
    {{with .Sections}}{{template "TOC" .}}{{end}}
  {{end}}
  </ul>
{{end}}

{{define "newline"}}
{{/* No automatic line break. Paragraphs are free-form. */}}
{{end}}

{/*
This is the action template.
It determines how the formatting actions are rendered.
*/}

{{define "section"}}
  <h{{len .Number}} id="TOC_{{.FormattedNumber}}">{{.FormattedNumber}} {{.Title}}</h{{len .Number}}>
  {{range .Elem}}{{elem $.Template .}}{{end}}
{{end}}

{{define "list"}}
  <ul>
  {{range .Bullet}}
    <li>{{style .}}</li>
  {{end}}
  </ul>
{{end}}

{{define "text"}}
  {{if .Pre}}
  <div class="code"><pre>{{range .Lines}}{{.}}{{end}}</pre></div>
  {{else}}
  <p>
    {{range $i, $l := .Lines}}{{if $i}}{{template "newline"}}
    {{end}}{{style $l}}{{end}}
  </p>
  {{end}}
{{end}}

{{define "code"}}
  <div class="code{{if playable .}} playground{{end}}" contenteditable="true" spellcheck="false">{{.Text}}</div>
{{end}}

{{define "image"}}
<div class="image">
  <img src="{{.URL}}"{{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}}>
</div>
{{end}}

{{define "iframe"}}
<iframe src="{{.URL}}"{{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}}></iframe>
{{end}}

{{define "link"}}<p class="link"><a href="{{.URL}}" target="_blank">{{style .Label}}</a></p>{{end}}

{{define "html"}}{{.HTML}}{{end}}

{{define "caption"}}<figcaption>{{style .Text}}</figcaption>{{end}}
`

	ARTICLE = `
Title
Sub Title
2 Jan 2015
Tags:test

secondarykey
Programer
secondarykey@gmail.com
http://github.com/secondarykey
@secondarykey

* Page

Description

* More Page
`
)
