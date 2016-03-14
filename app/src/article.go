package blog

import (
	"github.com/gorilla/mux"
	"golang.org/x/tools/present"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"

	"bufio"
	"bytes"
	"html/template"
	"net/http"
	"strings"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = createTemplate()
	if err != nil {
		panic(err)
	}
}

func createArticleHandler(w http.ResponseWriter, r *http.Request) {
	id, err := createArticle(r)
	if err != nil {
	}

	// Render Editor
	http.Redirect(w, r, "/admin/article/edit/"+id, 301)
}

func editArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["key"]

	c := appengine.NewContext(r)

	art, err := getArticle(r, name)
	if err != nil {
		log.Infof(c, err.Error())
		//NOT FOUND
		return
	}

	u, err := getUser(r)
	if err != nil {
		log.Infof(c, err.Error())
		//NOT FOUND
		return
	}

	s := struct {
		Article  *Article
		User     *User
		Markdown string
	}{art, u, string(art.Markdown)}

	adminRender(w, "./templates/admin/edit.tmpl", s)
}

func saveArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	c := appengine.NewContext(r)

	//http.Redirect(w, r, "/admin/article/edit/"+id, 301)
	_, err := updateArticle(r, id)
	if err != nil {
		log.Infof(c, err.Error())
		return
	}
	return
}

func publishArticleHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["key"]

	//http.Redirect(w, r, "/admin/article/edit/"+id, 301)
	err := updateHtml(r, id)
	if err != nil {
		c := appengine.NewContext(r)
		log.Infof(c, err.Error())
		return
	}

	return
}

func createTemplate() (*template.Template, error) {
	action := "templates/entry/action.tmpl"
	entry := "templates/entry/entry.tmpl"

	tmpl = present.Template()
	tmpl = tmpl.Funcs(template.FuncMap{"playable": playable})
	_, err := tmpl.ParseFiles(action, entry)
	if err != nil {
		return nil, err
	}
	return tmpl, nil

}

func createHtml(r *http.Request, art *Article) (datastore.ByteString, error) {

	//select user
	u, err := getUser(r)
	if err != nil {
		return nil, err
	}

	//create header
	header := art.Title + "\n\n" +
		u.Name + "\n" +
		u.Job + "\n" +
		u.Email + "\n" +
		u.URL + "\n" +
		u.TwitterId + "\n"

	txt := header + "\n" + string(art.Markdown)

	ctx := present.Context{ReadFile: readFile}
	reader := strings.NewReader(txt)
	doc, err := ctx.Parse(reader, "blog.article", 0)
	if err != nil {
		return nil, err
	}

	rtn := struct {
		*present.Doc
		Template    *template.Template
		PlayEnabled bool
	}{doc, tmpl, true}

	//Render
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	err = tmpl.ExecuteTemplate(writer, "root", rtn)

	if err != nil {
		return nil, err
	}
	writer.Flush()

	return b.Bytes(), nil
}

func readFile(name string) ([]byte, error) {
	//select file data
	return nil, nil
}
