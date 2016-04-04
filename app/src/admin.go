package blog

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
	"html/template"
	"net/http"
)

func adminRender(w http.ResponseWriter, tName string, obj interface{}) {

	funcMap := template.FuncMap{"convert": convert}
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

	var u *User
	var err error
	if r.Method == "POST" {
		u, err = putUser(r)
	} else {
		u, err = getUser(r)
	}

	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}

	adminRender(w, "./templates/admin/profile.tmpl", u)
}

func uploadAvatarHandler(w http.ResponseWriter, r *http.Request) {
	err := saveAvatar(r)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/profile", 301)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	u := user.Current(c)
	if u == nil {
		url, err := user.LoginURL(c, "/admin/")
		if err != nil {
			errorPage(w, "InternalServerError", err.Error(), 500)
			return
		}
		http.Redirect(w, r, url, 301)
		return
	}

	//exist user

	articles, err := selectArticle(r, 0)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
	adminRender(w, "./templates/admin/top.tmpl", articles)
}
