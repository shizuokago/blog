package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

var indexTmpl *template.Template

func init() {

	funcMap := template.FuncMap{"convert": Convert}
	var err error
	indexTmpl, err = template.New("root").Funcs(funcMap).ParseFiles("./cmd/static/templates/index.tmpl")
	if err != nil {
		log.Println(err)
	}
}

func topHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("top")

	var err error
	vals := r.URL.Query()
	ps := vals["p"]
	p := 1

	if len(ps) > 0 {
		pbuf := ps[0]
		p, err = strconv.Atoi(pbuf)
		if err != nil {
			ErrorPage(w, "Bad Request", err, 400)
			return
		}
	}

	//ページデータ
	htmls, err := datastore.SelectHtml(r, p)
	if err != nil {
		ErrorPage(w, "Not Found", err, 404)
		return
	}

	next := p + 1
	prev := p - 1
	flag := true
	if prev <= 0 {
		flag = false
	}

	//ブログデータ
	bgd := datastore.GetBlog(r)
	data := struct {
		Blog  *datastore.Blog
		HTMLs []datastore.Html
		Next  string
		Prev  string
		PFlag bool
	}{bgd, htmls, strconv.Itoa(next), strconv.Itoa(prev), flag}

	err = indexTmpl.Execute(w, data)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}
}

func entryHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	id := strings.Replace(url, "/entry/", "", 1)

	data, err := datastore.GetHtmlData(r, id)
	if err != nil {
		ErrorPage(w, "Not Found", err, 404)
		return
	}

	log.Println(string(data.Content))

	_, err = w.Write(data.Content)
	if err != nil {
		log.Println(err.Error())
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	t := "Not Found"
	m := "Page is Not Found"
	code := http.StatusNotFound
	ErrorPage(w, t, fmt.Errorf(m), code)
}
