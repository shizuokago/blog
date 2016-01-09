package main

import (
	"html/template"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	_ "golang.org/x/tools/playground"
	"golang.org/x/tools/present"
)

var basePath = "./"
var articleTemplate *template.Template

func init() {
	initTemplates(basePath)
	present.PlayEnabled = true

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")
	mime.AddExtensionType(".article", "text/html")

	r := mux.NewRouter()

	r.HandleFunc("/", topHandler).Methods("GET")
	r.HandleFunc("/entry/{key}", entryHandler).Methods("GET")
	r.HandleFunc("/engine/", engineHandler).Methods("GET")
	r.HandleFunc("/admin/profile", profileHandler)
	r.HandleFunc("/admin/", adminHandler).Methods("GET")

	r.HandleFunc("/admin/article/create", createArticleHandler).Methods("POST")

	//r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	http.Handle("/", r)
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}

func initTemplates(base string) error {
	// Locate the template file.
	actionTmpl := filepath.Join(base, "templates/action.tmpl")
	contentTmpl := filepath.Join(base, "templates/article.tmpl")

	// Read and parse the input.
	tmpl := present.Template()
	tmpl = tmpl.Funcs(template.FuncMap{"playable": playable})
	if _, err := tmpl.ParseFiles(actionTmpl, contentTmpl); err != nil {
		return err
	}
	articleTemplate = tmpl
	return nil
}

func renderDoc(w io.Writer, docFile string) error {
	doc, err := parse(docFile, 0)
	if err != nil {
		return err
	}
	return doc.Render(w, articleTemplate)
}

func parse(name string, mode present.ParseMode) (*present.Doc, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return present.Parse(f, name, 0)
}
