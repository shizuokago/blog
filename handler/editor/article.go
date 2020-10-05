package editor

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
	. "github.com/shizuokago/blog/handler/internal/form"
	"github.com/shizuokago/blog/logic"
)

func createArticleHandler(w http.ResponseWriter, r *http.Request) {

	file, data, err := GetFile(r)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}

	ctx := r.Context()

	blog := datastore.GetBlog(ctx)

	article := CreateNewArticle(blog)

	id, err := datastore.CreateArticle(ctx, article, file, data)
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

	ctx := r.Context()
	art, err := datastore.GetArticle(ctx, name)
	if err != nil {
		ErrorPage(w, "Key Error", err, 400)
		return
	}

	u, err := datastore.GetUser(ctx, user.Email)
	if err != nil {
		ErrorPage(w, "User Error", err, 401)
		return
	}

	auto := u.AutoSave
	autosave := ""
	if auto {
		autosave = "on"
	}

	bgd := datastore.GetBlog(ctx)
	s := struct {
		Article  *datastore.Article
		User     *datastore.User
		Markdown string
		BlogName string
		AutoSave string
	}{art, u, string(art.Markdown), bgd.Name, autosave}

	adminRender(w, "./cmd/templates/admin/edit.tmpl", s)
}

func saveArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	ctx := r.Context()

	art, err := CreateArticle(r)
	if err != nil {
		ErrorPage(w, "Argument Error", err, 400)
		return
	}

	err = datastore.UpdateArticle(ctx, id, art)
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

	art, err := CreateArticle(r)
	if err != nil {
		ErrorPage(w, "Argument Error", err, 400)
		return
	}

	ctx := r.Context()

	u, err := datastore.GetUser(ctx, user.Email)
	if err != nil {
		ErrorPage(w, "Get User Error", err, 400)
		return
	}

	p, err := logic.CreateHTML(ctx, id, art, u)
	if err != nil {
		ErrorPage(w, "Create HTML Error", err, 400)
		return
	}

	err = datastore.UpdateHTML(ctx, id, p, art)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	return
}

func deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["key"]

	ctx := r.Context()

	err := datastore.DeleteArticle(ctx, id)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/", 301)
}

func privateArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	ctx := r.Context()

	err := datastore.DeleteHTML(ctx, id)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/", 301)
}
