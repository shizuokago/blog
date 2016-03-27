package blog

import (
	"github.com/gorilla/mux"

	"html/template"
	"log"
	"net/http"
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

	curS := vals["cursor"]
	prevS := vals["prev"]

	cursor := ""
	if len(curS) > 0 {
		cursor = curS[0]
	}
	prev := ""
	if len(prevS) > 0 {
		prev = prevS[0]
	}

	htmls, next, err := selectHtml(r, cursor)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	if next == cursor {
		cursor = prev
	}

	data := struct {
		BlogName string
		HTMLs    []Html
		Next     string
		Prev     string
		Cursor   string
	}{blog.Name, htmls, next, prev, cursor}

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
