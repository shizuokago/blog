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

	http.Handle("/js/", fs)
	http.Handle("/css/", fs)
	http.Handle("/images/", fs)
	http.Handle("/favicon.ico", fs)

	//Deprecated
	//static := http.FileServer(http.Dir("./cmd"))
	//http.Handle("/static", static)

	//既存のブログが見ている
	http.Handle("/static/css/", fs)

	return nil
}
