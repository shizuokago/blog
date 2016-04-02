package blog

import (
	"github.com/gorilla/mux"
	"net/http"
)

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := createArticle(r)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
	}
	http.Redirect(w, r, "/admin/article/edit/"+id, 301)
}

func editArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["key"]

	art, err := getArticle(r, name)
	if err != nil {
		errorPage(w, "Key Error", err.Error(), 400)
		return
	}

	u, err := getUser(r)
	if err != nil {
		errorPage(w, "User Error", err.Error(), 401)
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
	_, err := updateArticle(r, id)
	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
	}
	return
}

func publishArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	err := updateHtml(r, id)
	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
	}

	return
}

func deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["key"]
	err := deleteArticle(r, id)
	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
	}
	http.Redirect(w, r, "/admin/", 301)
}
