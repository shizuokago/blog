package internal

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

var errorTmpl *template.Template

func init() {
	var err error
	errorTmpl, err = template.New("root").ParseFiles("./cmd/templates/error.tmpl")
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

	err = errorTmpl.Execute(w, data)
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
