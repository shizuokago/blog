package blog

import (
	"github.com/gorilla/mux"

	"html/template"
	"net/http"
	"strconv"
)

var indexTmpl *template.Template
var errorTmpl *template.Template

func init() {

	funcMap := template.FuncMap{"convert": convert}

	var err error
	indexTmpl, err = template.New("root").Funcs(funcMap).ParseFiles("./templates/index.tmpl")
	if err != nil {
		panic(err)
	}

	errorTmpl, err = template.New("root").ParseFiles("./templates/error.tmpl")
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

	htmls, err := selectHtml(r, p)
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

	data := struct {
		Blog  Blog
		HTMLs []Html
		Next  string
		Prev  string
		PFlag bool
	}{blog, htmls, strconv.Itoa(next), strconv.Itoa(prev), flag}

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
	data, err := getHtmlData(r, id)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 404)
		return
	}

	w.Write(data.Content)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	t := "Not Found"
	m := "Page is Not Found"
	code := http.StatusNotFound
	errorPage(w, t, m, code)
}
