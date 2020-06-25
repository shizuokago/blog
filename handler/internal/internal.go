package internal

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

var errorTmpl *template.Template

func init() {
	var err error
	errorTmpl, err = template.New("root").ParseFiles("./cmd/static/templates/error.tmpl")
	if err != nil {
		panic(err)
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

func ErrorPage(w http.ResponseWriter, t, m string, code int) {

	log.Println(t)
	log.Println(m)
	log.Println(code)

	data := struct {
		Code    int
		Title   string
		Message string
	}{code, t, m}

	w.WriteHeader(data.Code)

	err := errorTmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func ErrorJson(w http.ResponseWriter, t, m string, code int) {

	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	d := map[string]interface{}{
		"success": false,
		"title":   t,
		"message": m,
	}
	enc.Encode(d)
}
