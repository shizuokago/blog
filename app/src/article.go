package blog

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

	art, err := getArticle(r, name)
	if err != nil {
		log.Infof(c, err.Error())
		//NOT FOUND
		return
	}

	u, err := getUser(r)
	if err != nil {
		log.Infof(c, err.Error())
		//NOT FOUND
		return
	}

	s := struct {
		Article  *Article
		User     *User
		Markdown string
		BlogName string
	}{art, u, string(art.Markdown), blog.Name}

	adminRender(w, "./templates/admin/edit.tmpl", s)
}

func saveArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	c := appengine.NewContext(r)

	//http.Redirect(w, r, "/admin/article/edit/"+id, 301)
	_, err := updateArticle(r, id)
	if err != nil {
		log.Infof(c, err.Error())
		return
	}
	return
}

func publishArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	//http.Redirect(w, r, "/admin/article/edit/"+id, 301)
	err := updateHtml(r, id)
	if err != nil {
		c := appengine.NewContext(r)
		log.Infof(c, err.Error())
		return
	}

	return
}
