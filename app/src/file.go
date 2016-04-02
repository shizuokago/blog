package blog

import (
	"net/http"
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

func fileViewHandler(w http.ResponseWriter, r *http.Request) {
	adminRender(w, "./templates/admin/file.tmpl", nil)
}

func fileUploadHandler(w http.ResponseWriter, r *http.Request) {

	err := saveFile(r, "", FILE_TYPE_DATA)
	if err != nil {
		errorPage(w, "InternalServerError", err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/admin/file/view", 301)
}
