package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func init() {
}

func GrantFS(f http.FileSystem, g string) *grantFS {
	var grant grantFS
	grant.fs = f
	grant.prefix = g
	return &grant
}

type grantFS struct {
	fs     http.FileSystem
	prefix string
}

func (g grantFS) Open(name string) (http.File, error) {
	n := g.prefix + name
	return g.fs.Open(n)
}

func WasmServer(f http.FileSystem) func(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(f)
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/wasm")
		w.Header().Set("Content-Encoding", "gzip")
		handler.ServeHTTP(w, r)
	}
}

func Convert(t time.Time) string {
	if t.IsZero() {
		return "None"
	}
	jst, _ := time.LoadLocation("Asia/Tokyo")
	jt := t.In(jst)
	return jt.Format("2006/01/02 15:04")
}

func ErrorPage(w http.ResponseWriter, t string, err error, code int) {

	buf := fmt.Sprintf("%+v", err)

	if code != 404 {
		log.Println(buf)
	}

	data := struct {
		Code    int
		Title   string
		Message string
		Detail  string
	}{code, t, err.Error(), buf}

	w.WriteHeader(data.Code)

	tmpl, err := GetTemplate(nil, "error.tmpl")
	if err != nil {
		log.Println("error page getTemplate", err.Error())
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("ErrorPage() write error:", err.Error())
	}
}

func ErrorJson(w http.ResponseWriter, t string, err error, code int) {

	buf := fmt.Sprintf("%+v", err)

	enc := json.NewEncoder(w)
	d := map[string]interface{}{
		"success": false,
		"title":   t,
		"message": err.Error(),
		"detail":  buf,
	}

	w.WriteHeader(code)
	err = enc.Encode(d)
	if err != nil {
		log.Println("ErrorJson() write error:", err.Error())
	}
}
