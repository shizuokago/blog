package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {

	AddCacheHeader(w, r)

	url := r.URL.Path
	name := strings.Replace(url, "/file/", "", 1)

	file, err := datastore.GetFileData(r, name)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}

	if file == nil {
		ErrorPage(w, "Not Found", fmt.Errorf("FileNotFound:"+name), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", file.Mime)

	//set MIME
	_, err = w.Write(file.Content)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
}
