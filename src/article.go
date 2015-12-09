package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/tools/present"
)

func init() {
	http.HandleFunc("/", topHandler)
	http.HandleFunc("/entry/", entryHandler)
	http.HandleFunc("/engine/", engineHandler)
}

func topHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := filepath.Join("./", "templates/index.tmpl")
	tmplObj, err := template.ParseFiles(tmpl)
	if err != nil {
	}
	tmplObj.Execute(w, nil)
}

func entryHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := filepath.Join("./", "templates/entry.tmpl")
	tmplObj, err := template.ParseFiles(tmpl)
	if err != nil {
	}
	tmplObj.Execute(w, nil)
}

// dirHandler serves a directory listing for the requested path, rooted at basePath.
func engineHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.Error(w, "not found", 404)
		return
	}
	const base = "."
	name := filepath.Join(base, r.URL.Path)
	err := renderDoc(w, name)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 500)
	}
}

var (
	contentTemplate map[string]*template.Template
)

func initTemplates(base string) error {
	// Locate the template file.
	actionTmpl := filepath.Join(base, "templates/action.tmpl")

	contentTemplate = make(map[string]*template.Template)

	for ext, contentTmpl := range map[string]string{
		".article": "article.tmpl",
	} {
		contentTmpl = filepath.Join(base, "templates", contentTmpl)

		// Read and parse the input.
		tmpl := present.Template()
		tmpl = tmpl.Funcs(template.FuncMap{"playable": playable})
		if _, err := tmpl.ParseFiles(actionTmpl, contentTmpl); err != nil {
			return err
		}
		contentTemplate[ext] = tmpl
	}
	return nil
}

// renderDoc reads the present file, gets its template representation,
// and executes the template, sending output to w.
func renderDoc(w io.Writer, docFile string) error {
	// Read the input and build the doc structure.
	doc, err := parse(docFile, 0)
	if err != nil {
		return err
	}

	// Find which template should be executed.
	tmpl := contentTemplate[filepath.Ext(docFile)]

	// Execute the template.
	return doc.Render(w, tmpl)
}

func parse(name string, mode present.ParseMode) (*present.Doc, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return present.Parse(f, name, 0)
}
