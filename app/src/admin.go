package blog

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
	"google.golang.org/appengine/user"

	"html/template"
	"net/http"
	"strconv"
)

func adminRender(w http.ResponseWriter, tName string, obj interface{}) {

	funcMap := template.FuncMap{"convert": convert, "deleteDir": deleteDir}
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
	/*
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
	*/

	vals := r.URL.Query()
	ps := vals["p"]
	p := 1
	if len(ps) > 0 {
		pbuf := ps[0]
		p, err := strconv.Atoi(p)
		if err != nil {
			errorPage(w, "Bad Request", err.Error(), 400)
			return
		}
	}

	articles, err := selectArticle(r, p)
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
		Articles []Article
		Next     string
		Prev     string
		PFlag    bool
	}{articles, strconv.Itoa(next), strconv.Itoa(prev), flag}

	adminRender(w, "./templates/admin/top.tmpl", data)
}
