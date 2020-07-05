package editor

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

func deleteDir(s string) string {
	ds := []byte(s)
	return string(ds[5:])
}

func adminRender(w http.ResponseWriter, tName string, obj interface{}) {

	funcMap := template.FuncMap{"convert": Convert, "deleteDir": deleteDir}
	tmpl, err := template.New("root").Funcs(funcMap).ParseFiles("./cmd/templates/admin/layout.tmpl", tName)
	if err != nil {
		ErrorPage(w, "Template Parse Error", err, 500)
		return
	}

	err = tmpl.Execute(w, obj)
	if err != nil {
		ErrorPage(w, "Template Execute Error", err, 500)
		return
	}
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
			ErrorPage(w, "Bad Request", err, 400)
			return
		}
	}

	articles, err := datastore.SelectArticle(r, p)
	if err != nil {
		ErrorPage(w, "Not Found", err, 404)
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

	adminRender(w, "./cmd/templates/admin/top.tmpl", data)
}
