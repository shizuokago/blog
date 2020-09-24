package handler

import (
	"mime"
	"net/http"
)

func registerStatic() error {

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")

	//r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	fs := http.FileServer(http.Dir("./cmd/assets"))

	http.Handle("/favicon.ico", fs)
	http.Handle("/js/", fs)
	http.Handle("/images/", fs)

	http.HandleFunc("/css/styles.css", deprecatedStyles)
	http.HandleFunc("/static/css/styles.css", deprecatedStyles)

	//Deprecated
	//static := http.FileServer(http.Dir("./cmd"))
	//http.Handle("/static", static)
	//既存のブログが見ている

	return nil
}

func deprecatedStyles(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/file/data/styles.css", 302)
}
