package editor

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
	. "github.com/shizuokago/blog/handler/internal/form"
)

func existsFileHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	id := r.FormValue("fileName")

	ctx := r.Context()

	flag, err := datastore.ExistsFile(ctx, id, datastore.FILE_TYPE_DATA)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}

	enc := json.NewEncoder(w)
	d := map[string]bool{"exists": flag}
	enc.Encode(d)

	return
}

func viewFileHandler(w http.ResponseWriter, r *http.Request) {

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

	ctx := r.Context()
	files, err := datastore.SelectFile(ctx, p)
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

	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	data := struct {
		Files []datastore.File
		Next  string
		Prev  string
		PFlag bool
	}{files, strconv.Itoa(next), strconv.Itoa(prev), flag}
	adminRender(w, "file.tmpl", data)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	f, d, err := GetFile(r)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}

	id := r.FormValue("FileName")
	p := datastore.FileParam{
		File:     f,
		FileData: d,
	}

	ctx := r.Context()
	err = datastore.SaveFile(ctx, id, datastore.FILE_TYPE_DATA, &p)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}

	http.Redirect(w, r, "/admin/file/view", 302)
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("FileName")

	ctx := r.Context()
	err := datastore.DeleteFile(ctx, id)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/file/view", 302)
}

func saveBackgroundHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["key"]

	f, d, err := GetFile(r)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	p := datastore.FileParam{
		File:     f,
		FileData: d,
	}

	ctx := r.Context()
	err = datastore.SaveBackgroundImage(ctx, name, &p)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/article/edit/"+name, 302)
}

func deleteBackgroundHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["key"]

	ctx := r.Context()
	err := datastore.DeleteBackgroundImage(ctx, name)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/article/edit/"+name, 302)
}
