package handler

import (
	"net/http"
	"strings"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	name := strings.Replace(url, "/file/", "", 1)

	file, err := datastore.GetFileData(r, name)
	if err != nil {
		ErrorPage(w, "InternalServerError", err.Error(), 500)
		return
	}

	if file == nil {
		ErrorPage(w, "Not Found", "File Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", file.Mime)

	//set MIME
	_, err = w.Write(file.Content)
	if err != nil {
		ErrorPage(w, "InternalServerError", err.Error(), 500)
		return
	}
}
