package editor

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := datastore.CreateArticle(r)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/article/edit/"+id, 301)
}

func editArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["key"]

	user, err := GetSession(r)

	art, err := datastore.GetArticle(r, name)
	if err != nil {
		ErrorPage(w, "Key Error", err, 400)
		return
	}

	u, err := datastore.GetUser(r, user.Email)
	if err != nil {
		ErrorPage(w, "User Error", err, 401)
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

	adminRender(w, "./cmd/static/templates/admin/edit.tmpl", s)
}

func saveArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]
	_, err := datastore.UpdateArticle(r, id, time.Time{})
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}
	return
}

func publishArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	user, err := GetSession(r)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	err = datastore.UpdateHtml(r, user.Email, id)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	return
}

func deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["key"]
	err := datastore.DeleteArticle(r, id)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/", 301)
}

func privateArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]
	err := datastore.DeleteHtml(r, id)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/", 301)
}
