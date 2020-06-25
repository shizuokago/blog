package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/shizuokago/blog/datastore"
)

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := datastore.CreateArticle(r)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
	}
	http.Redirect(w, r, "/admin/article/edit/"+id, 301)
}

func editArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["key"]

	art, err := datastore.GetArticle(r, name)

	if err != nil {
		errorPage(w, "Key Error", err.Error(), 400)
		return
	}

	u, err := datastore.GetUser(r)
	if err != nil {
		errorPage(w, "User Error", err.Error(), 401)
		return
	}

	auto := u.AutoSave
	autosave := ""
	if auto {
		autosave = "on"
	}

	bgd := datastore.GetBlog(r)
	s := struct {
		Article  *datastore.Article
		User     *datastore.User
		Markdown string
		BlogName string
		AutoSave string
	}{art, u, string(art.Markdown), bgd.Name, autosave}

	adminRender(w, "./templates/admin/edit.tmpl", s)
}

func saveArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]
	_, err := datastore.UpdateArticle(r, id, time.Time{})
	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
	}
	return
}

func publishArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	err := datastore.UpdateHtml(r, id)
	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
	}

	return
}

func deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["key"]
	err := datastore.DeleteArticle(r, id)
	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
	}
	http.Redirect(w, r, "/admin/", 301)
}

func privateArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]
	err := datastore.DeleteHtml(r, id)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
	}
	http.Redirect(w, r, "/admin/", 301)
}