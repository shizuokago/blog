package blog

import (
	"github.com/gorilla/mux"

	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var indexTmpl *template.Template

func init() {
	var err error
	tmpl := filepath.Join("./", "templates/index.tmpl")
	indexTmpl, err = template.ParseFiles(tmpl)
	if err != nil {
		log.Println(err)
	}
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	//Get PageNum
	page := 1
	//Get PageList
	htmls, err := selectHtml(r, page)
	if err != nil {
		log.Println(err)
	}

	data := struct {
		BlogName string
		HTMLs    []Html
	}{blog.Name, htmls}

	indexTmpl.Execute(w, data)
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
