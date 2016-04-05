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
	vals := r.URL.Query()

	ps := vals["p"]
	p := "1"
	if len(ps) > 0 {
		p = ps[0]
	}

	item, err := memcache.Get(c, "article_"+p+"_cursor")
	cursor := ""
	if err == nil {
		cursor = string(item.Value)
	}
	articles, nextC, err := selectArticle(r, cursor)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 404)
		return
	}

	t, err := strconv.Atoi(p)
	if err != nil {
		errorPage(w, "Page Error", err.Error(), 400)
		return
	}

	next := t + 1
	prev := t - 1
	flag := true
	if prev <= 0 {
		flag = false
	}

	err = memcache.Set(c, &memcache.Item{
		Key:   "article_" + strconv.Itoa(next) + "_cursor",
		Value: []byte(nextC),
	})

	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
	}

	data := struct {
		Articles []Article
		Next     string
		Prev     string
		PFlag    bool
	}{articles, strconv.Itoa(next), strconv.Itoa(prev), flag}

	//articles, err := selectArticle(r, 0)
	//if err != nil {
	//errorPage(w, "InternalServerError", err.Error(), 500)
	//return
	//}

	adminRender(w, "./templates/admin/top.tmpl", data)
}
