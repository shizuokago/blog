package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

func adminRender(w http.ResponseWriter, tName string, obj interface{}) {

	funcMap := template.FuncMap{"convert": Convert, "deleteDir": deleteDir}
	tmpl, err := template.New("root").Funcs(funcMap).ParseFiles("./templates/admin/layout.tmpl", tName)
	if err != nil {
		errorPage(w, "Template Parse Error", err.Error(), 500)
		return
	}

	err = tmpl.Execute(w, obj)
	if err != nil {
		errorPage(w, "Template Execute Error", err.Error(), 500)
		return
	}
}

func profileHandler(w http.ResponseWriter, r *http.Request) {

	var u *datastore.User
	var err error
	if r.Method == "POST" {
		u, err = datastore.PutInformation(r)
	} else {
		u, err = datastore.GetUser(r)
	}

	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}

	if u == nil {
		u = &datastore.User{}
	}

	bgd := datastore.GetBlog(r)
	data := struct {
		Blog *datastore.Blog
		User *datastore.User
	}{bgd, u}

	adminRender(w, "./templates/admin/profile.tmpl", data)
}

func uploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	err := datastore.SaveAvatar(r)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/profile", 301)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {

	bgd := datastore.GetBlog(r)
	if bgd.Name == "" {
		http.Redirect(w, r, "/admin/profile", 301)
		return
	}

	var err error
	vals := r.URL.Query()
	ps := vals["p"]
	p := 1

	if len(ps) > 0 {
		pbuf := ps[0]
		p, err = strconv.Atoi(pbuf)
		if err != nil {
			errorPage(w, "Bad Request", err.Error(), 400)
			return
		}
	}

	articles, err := datastore.SelectArticle(r, p)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 404)
		return
	}

	next := p + 1
	prev := p - 1
	flag := true
	if prev <= 0 {
		flag = false
	}

	data := struct {
		Articles []datastore.Article
		Next     string
		Prev     string
		PFlag    bool
	}{articles, strconv.Itoa(next), strconv.Itoa(prev), flag}

	adminRender(w, "./templates/admin/top.tmpl", data)
}
