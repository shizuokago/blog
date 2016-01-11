package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func init() {

}

func topHandler(w http.ResponseWriter, r *http.Request) {

	//Get PageNum

	//Get PageList
	articles, err := selectArticle(r, 0)
	if err != nil {
	}

	tmpl := filepath.Join("./", "templates/index.tmpl")
	tmplObj, err := template.ParseFiles(tmpl)
	if err != nil {
	}
	tmplObj.Execute(w, articles)
}

func entryHandler(w http.ResponseWriter, r *http.Request) {

	//Get Key
	//vars := mux.Vars(r)

	//Find Article

	//Render

	tmpl := filepath.Join("./", "templates/entry.tmpl")
	tmplObj, err := template.ParseFiles(tmpl)
	if err != nil {
	}
	tmplObj.Execute(w, nil)
}

func engineHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.Error(w, "not found", 404)
		return
	}
	const base = "."
	name := filepath.Join(base, r.URL.Path)
	err := renderDoc(w, name)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
	}
}
