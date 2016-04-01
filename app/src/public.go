package blog

import (
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"

	"html/template"
	"log"
	"net/http"
	"strconv"
)

var indexTmpl *template.Template

func init() {

	funcMap := template.FuncMap{"convert": convert}

	var err error
	indexTmpl, err = template.New("root").Funcs(funcMap).ParseFiles("./templates/index.tmpl")
	//indexTmpl, err = template.ParseFiles("./templates/index.tmpl")
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
	item, err := memcache.Get(c, "paging_"+p+"_cursor")
	cursor := ""
	if err == nil {
		cursor = string(item.Value)
	}

	htmls, nextC, err := selectHtml(r, cursor)
	if err != nil {
		panic(err)
	}

	t, _ := strconv.Atoi(p)
	next := t + 1
	prev := t - 1
	flag := true
	if prev <= 0 {
		flag = false
	}

	memcache.Set(c, &memcache.Item{
		Key:   "paging_" + strconv.Itoa(next) + "_cursor",
		Value: []byte(nextC),
	})

	data := struct {
		Blog  Blog
		HTMLs []Html
		Next  string
		Prev  string
		PFlag bool
	}{blog, htmls, strconv.Itoa(next), strconv.Itoa(prev), flag}

	err = indexTmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func entryHandler(w http.ResponseWriter, r *http.Request) {
	//Get Key
	vars := mux.Vars(r)
	id := vars["key"]

	data, err := getHtmlData(r, id)
	if err != nil {
		log.Println(err)
	}

	w.Write(data.Content)
}
