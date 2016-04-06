package blog

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"

	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	name := strings.Replace(url, "/file/", "", 1)

	file, err := getFileData(r, name)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), http.StatusNotFound)
		return
	}

	//set MIME
	_, err = w.Write(file.Content)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 500)
		return
	}
}

func existsFileHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	id := r.FormValue("fileName")
	flag, err := existsFile(r, id, FILE_TYPE_DATA)
	if err != nil {
		errorJson(w, "InternalServerError", err.Error(), 500)
		return
	}

	enc := json.NewEncoder(w.Body)
	d := map[string]bool{"exists": flag}
	enc.Encode(d)

	return
}

func viewFileHandler(w http.ResponseWriter, r *http.Request) {

	vals := r.URL.Query()

	ps := vals["p"]
	p := "1"
	if len(ps) > 0 {
		p = ps[0]
	}

	c := appengine.NewContext(r)
	item, err := memcache.Get(c, "file_"+p+"_cursor")
	cursor := ""
	if err == nil {
		cursor = string(item.Value)
	}

	files, nextC, err := selectFile(r, cursor)
	if err != nil {
		errorPage(w, "Not Found", err.Error(), 404)
		return
	}

	t, err := strconv.Atoi(p)
	if err != nil {
		errorPage(w, "Page Error", err.Error(), 400)
		return
	}

	next := t + 1
	prev := t - 1
	flag := true
	if prev <= 0 {
		flag = false
	}

	err = memcache.Set(c, &memcache.Item{
		Key:   "file_" + strconv.Itoa(next) + "_cursor",
		Value: []byte(nextC),
	})

	if err != nil {
		errorPage(w, "Internal Server Error", err.Error(), 500)
		return
	}
	data := struct {
		Files []File
		Next  string
		Prev  string
		PFlag bool
	}{files, strconv.Itoa(next), strconv.Itoa(prev), flag}
	adminRender(w, "./templates/admin/file.tmpl", data)
}

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	err := saveFile(r, r.FormValue("FileName"), FILE_TYPE_DATA)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/admin/file/view", 301)
}

func deleteFileHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("FileName")
	err := deleteFile(r, "data/"+filename)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/admin/file/view", 301)
}
