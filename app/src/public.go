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

	err = indexTmpl.Execute(w, data)
	if err != nil {
		log.Println(indexTmpl)
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
