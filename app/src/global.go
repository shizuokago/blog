package blog

import (
	"bufio"
	"bytes"
	"encoding/json"
	"html/template"
	"image"
	"image/jpeg"
	"io/ioutil"
	"mime"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"

	_ "golang.org/x/tools/playground"
	"golang.org/x/tools/present"
)

type Blog struct {
	Name   string
	Author string
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

	r.HandleFunc("/", topHandler).Methods("GET")
	r.HandleFunc("/entry/{key}", entryHandler).Methods("GET")

	r.HandleFunc("/admin/profile/upload", uploadAvatarHandler).Methods("POST")
	r.HandleFunc("/admin/profile", profileHandler)
	r.HandleFunc("/admin/", adminHandler).Methods("GET")

	r.HandleFunc("/admin/article/create", createArticleHandler).Methods("POST")
	r.HandleFunc("/admin/article/edit/{key}", editArticleHandler).Methods("GET")

	r.HandleFunc("/admin/article/save/{key}", saveArticleHandler).Methods("POST")
	r.HandleFunc("/admin/article/publish/{key}", publishArticleHandler).Methods("POST")

	r.HandleFunc("/admin/article/delete/{key}", deleteArticleHandler).Methods("GET")

	//r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
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
	jst, _ := time.LoadLocation("Asia/Tokyo")
	jt := t.In(jst)
	return jt.Format("2006/01/02 15:04")
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

func readFile(name string) ([]byte, error) {
	//select file data
	return nil, nil
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
		AuthorID    string
		StringID    string
		BlogName    string
		HTML        *Html
	}{doc, tmpl, true, u.Key.StringID(), art.Key.StringID(), blog.Name, html}

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

func resizeImage(b []byte) ([]byte, error) {

	buff := bytes.NewBuffer(b)

	img, _, err := image.Decode(buff)
	if err != nil {
		return nil, err
	}

	m := resize.Resize(1000, 0, img, resize.Lanczos3)
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, m, nil); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
