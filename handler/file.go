package handler

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"github.com/shizuokago/blog/datastore"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	name := strings.Replace(url, "/file/", "", 1)

	file, err := datastore.GetFileData(r, name)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}

	if file == nil {
		errorPage(w, "Not Found", "File Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", file.Mime)

	//set MIME
	_, err = w.Write(file.Content)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
}

func existsFileHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	id := r.FormValue("fileName")

	c := appengine.NewContext(r)
	log.Infof(c, id)

	flag, err := datastore.ExistsFile(r, id, datastore.FILE_TYPE_DATA)
	if err != nil {
		errorJson(w, "InternalServerError", err.Error(), 500)
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
			errorPage(w, "Bad Request", err.Error(), 400)
			return
		}
	}

	files, err := datastore.SelectFile(r, p)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 404)
		return
	}

	next := p + 1
	prev := p - 1
	flag := true
	if prev <= 0 {
		flag = false
	}

	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
	}

	data := struct {
		Files []datastore.File
		Next  string
		Prev  string
		PFlag bool
	}{files, strconv.Itoa(next), strconv.Itoa(prev), flag}
	adminRender(w, "./templates/admin/file.tmpl", data)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	err := datastore.SaveFile(r, r.FormValue("FileName"), datastore.FILE_TYPE_DATA)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/admin/file/view", 301)
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("FileName")
	err := datastore.DeleteFile(r, "data/"+filename)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/file/view", 301)
}

func saveBackgroundHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["key"]

	err := datastore.SaveBackgroundImage(r, name)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/article/edit/"+name, 301)
}

func deleteBackgroundHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["key"]

	err := datastore.DeleteBackgroundImage(r, name)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/article/edit/"+name, 301)
}
