package main

import (
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"net/http"
)

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := createArticle(r)
	if err != nil {
	}
	// Render Editor
	http.Redirect(w, r, "/admin/article/edit/"+id, 301)
}

func editArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["key"]

	c := appengine.NewContext(r)
	log.Infof(c, name)

	art, err := getArticle(r, name)
	if err != nil {
		c := appengine.NewContext(r)
		log.Infof(c, err.Error())
	}
	adminRender(w, "./templates/admin/edit.tmpl", art)
}
