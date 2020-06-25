package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

var indexTmpl *template.Template
var errorTmpl *template.Template

func init() {

	funcMap := template.FuncMap{"convert": Convert}

	var err error
	indexTmpl, err = template.New("root").Funcs(funcMap).ParseFiles("./cmd/static/templates/index.tmpl")
	if err != nil {
		panic(err)
	}

	errorTmpl, err = template.New("root").ParseFiles("./cmd/static/templates/error.tmpl")
	if err != nil {
		panic(err)
	}
}

func topHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	vals := r.URL.Query()
	ps := vals["p"]
	p := 1

	if len(ps) > 0 {
		pbuf := ps[0]
		p, err = strconv.Atoi(pbuf)
		if err != nil {
			errorPage(w, "Bad Request", err.Error(), 400)
			return
		}
	}

	htmls, err := datastore.SelectHtml(r, p)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 404)
		return
	}

	next := p + 1
	prev := p - 1
	flag := true
	if prev <= 0 {
		flag = false
	}

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
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
	}
}

func entryHandler(w http.ResponseWriter, r *http.Request) {
	//Get Key
	vars := mux.Vars(r)
	id := vars["key"]
	data, err := datastore.GetHtmlData(r, id)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 404)
		return
	}

	w.Write([]byte(data.Content))
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	t := "Not Found"
	m := "Page is Not Found"
	code := http.StatusNotFound
	errorPage(w, t, m, code)
}
