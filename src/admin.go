package blog

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/user"
	"html/template"
	"net/http"
)

func adminRender(w http.ResponseWriter, tName string, obj interface{}) {

	tmpl, err := template.ParseFiles("./templates/admin/layout.tmpl", tName)
	if err != nil {

		//error page

		return
	}
	err = tmpl.Execute(w, obj)
	if err != nil {

		//error page

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

	// add error handling
	c := appengine.NewContext(r)
	if err != nil {
		log.Infof(c, "%T,%s", err, err.Error())
	}

	adminRender(w, "./templates/admin/profile.tmpl", u)
}

func adminHandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)

	log.Infof(c, "Test")
	u := user.Current(c)

	if u == nil {
		url, err := user.LoginURL(c, "/admin/")
		if err != nil {
		}

		http.Redirect(w, r, url, 301)
		return
	}

	//exists user

	log.Infof(c, "Article")
	//find article
	articles, err := selectArticle(r, 0)
	if err != nil {
	}

	log.Infof(c, "Render")
	adminRender(w, "./templates/admin/top.tmpl", articles)
}
