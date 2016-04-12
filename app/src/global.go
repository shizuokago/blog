package blog

import (
	"bufio"
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	_ "golang.org/x/tools/playground"
	"golang.org/x/tools/present"
)

type Blog struct {
	Name        string
	Author      string
	Tags        string
	Description string
	Template    string
}

var tmpl *template.Template
var blog = Blog{}

func init() {

	var err error
	tmpl, err = createTemplate()
	if err != nil {
		panic(err)
	}
	present.PlayEnabled = true

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")

	r := mux.NewRouter()

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	r.HandleFunc("/", topHandler).Methods("GET")
	r.HandleFunc("/entry/{key}", entryHandler).Methods("GET")

	r.HandleFunc("/admin/profile/upload", uploadAvatarHandler).Methods("POST")
	r.HandleFunc("/admin/profile", profileHandler)
	r.HandleFunc("/admin/", adminHandler).Methods("GET")

	r.HandleFunc("/admin/article/create", createArticleHandler).Methods("POST")
	r.HandleFunc("/admin/article/edit/{key}", editArticleHandler).Methods("GET")

	r.HandleFunc("/admin/article/save/{key}", saveArticleHandler).Methods("POST")
	r.HandleFunc("/admin/article/publish/{key}", publishArticleHandler).Methods("POST")
	r.HandleFunc("/admin/article/private/{key}", privateArticleHandler)

	r.HandleFunc("/admin/article/delete/{key}", deleteArticleHandler).Methods("GET")

	r.HandleFunc("/admin/article/bg/save/{key}", saveBackgroundHandler)
	r.HandleFunc("/admin/article/bg/delete/{key}", deleteBackgroundHandler)

	r.HandleFunc("/admin/file/view", viewFileHandler).Methods("GET")
	r.HandleFunc("/admin/file/upload", uploadFileHandler).Methods("POST")
	r.HandleFunc("/admin/file/delete", deleteFileHandler).Methods("POST")
	r.HandleFunc("/admin/file/exists", existsFileHandler).Methods("POST")

	http.HandleFunc("/file/", fileHandler)
	http.Handle("/", r)

	jsonString, err := ioutil.ReadFile("./static/blog.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonString, &blog)
	if err != nil {
		panic(err)
	}
}

func playable(c present.Code) bool {
	return present.PlayEnabled && c.Play && c.Ext == ".go"
}

func convert(t time.Time) string {
	if t.IsZero() {
		return "None"
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	jt := t.In(jst)
	return jt.Format("2006/01/02 15:04")
}

func deleteDir(s string) string {
	ds := []byte(s)
	return string(ds[5:])
}

func createTemplate() (*template.Template, error) {

	action := "templates/entry/action.tmpl"
	entry := "templates/entry/entry.tmpl"

	tmpl = present.Template()
	funcMap := template.FuncMap{
		"playable": playable,
		"convert":  convert,
	}
	tmpl = tmpl.Funcs(funcMap)
	_, err := tmpl.ParseFiles(action, entry)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

type FileDs struct {
	request *http.Request
}

func (ds FileDs) readFile(name string) ([]byte, error) {
	key := "data/" + name
	file, err := getFileData(ds.request, key)
	if err != nil {
	}
	return file.Content, nil
}

func createHtml(r *http.Request, art *Article, u *User, html *Html) ([]byte, error) {

	//create header
	header := art.Title + "\n\n" +
		u.Name + "\n" +
		u.Job + "\n" +
		u.Email + "\n" +
		u.URL + "\n" +
		u.TwitterId + "\n"

	txt := header + "\n" + string(art.Markdown)

	ds := FileDs{
		request: r,
	}

	ctx := present.Context{ReadFile: ds.readFile}

	reader := strings.NewReader(txt)
	doc, err := ctx.Parse(reader, "blog.article", 0)
	if err != nil {
		return nil, err
	}

	bd := Blog{
		Name:        blog.Name,
		Author:      html.Author,
		Tags:        art.Tags,
		Description: art.SubTitle,
	}

	rtn := struct {
		*present.Doc
		Template    *template.Template
		PlayEnabled bool
		AuthorID    string
		StringID    string
		Blog        Blog
		HTML        *Html
	}{doc, tmpl, true, u.Key.StringID(), art.Key.StringID(), bd, html}

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

func errorPage(w http.ResponseWriter, t, m string, code int) {
	data := struct {
		Blog    Blog
		Code    int
		Title   string
		Message string
	}{blog, code, t, m}

	w.WriteHeader(data.Code)
	err := errorTmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func errorJson(w http.ResponseWriter, t, m string, code int) {

	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	d := map[string]interface{}{
		"success": false,
		"title":   t,
		"message": m,
	}
	enc.Encode(d)
}
