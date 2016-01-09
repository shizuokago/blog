package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

func init() {

	r := mux.NewRouter()

	r.HandleFunc("/", topHandler).Methods("GET")
	r.HandleFunc("/entry/{key}", entryHandler).Methods("GET")
	r.HandleFunc("/engine/", engineHandler).Methods("GET")
	r.HandleFunc("/admin/profile", profileHandler)
	r.HandleFunc("/admin/", adminHandler).Methods("GET")

	//r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	http.Handle("/", r)
}

func topHandler(w http.ResponseWriter, r *http.Request) {

	//Get PageNum

	//Get PageList

	tmpl := filepath.Join("./", "templates/index.tmpl")
	tmplObj, err := template.ParseFiles(tmpl)
	if err != nil {
	}
	tmplObj.Execute(w, nil)
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
