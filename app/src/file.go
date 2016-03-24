package blog

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"net/http"
	"strings"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {

	url := r.URL.Path
	name := strings.Replace(url, "/file/", "", 1)

	c := appengine.NewContext(r)
	log.Infof(c, name)

	file, err := getFileData(r, name)
	if err != nil {
	}
	//set MIME
	w.Write(file.Content)
}
