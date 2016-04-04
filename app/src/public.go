package blog

import (
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"

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
	//indexTmpl, err = template.ParseFiles("./templates/index.tmpl")
	if err != nil {
		panic(err)
	}

	errorTmpl, err = template.New("root").ParseFiles("./templates/error.tmpl")
	if err != nil {
		panic(err)
	}
}

func topHandler(w http.ResponseWriter, r *http.Request) {

	vals := r.URL.Query()

	ps := vals["p"]
	p := "1"
	if len(ps) > 0 {
		p = ps[0]
	}

	c := appengine.NewContext(r)
	item, err := memcache.Get(c, "html_"+p+"_cursor")
	cursor := ""
	if err == nil {
		cursor = string(item.Value)
	}

	htmls, nextC, err := selectHtml(r, cursor)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 404)
		return
	}

	t, err := strconv.Atoi(p)
	if err != nil {
		errorPage(w, "Page Error", err.Error(), 400)
		return
	}

	next := t + 1
	prev := t - 1
	flag := true
	if prev <= 0 {
		flag = false
	}

	err = memcache.Set(c, &memcache.Item{
		Key:   "html_" + strconv.Itoa(next) + "_cursor",
		Value: []byte(nextC),
	})

	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
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
