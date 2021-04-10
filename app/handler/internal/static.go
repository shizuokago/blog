package internal

import (
	"mime"
	"net/http"
)

func init() {
}

func RegisterStatic() error {

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")

	//r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	fs := http.FileServer(GrantFS(statikFS, "/static"))

	http.Handle("/favicon.ico", fs)
	http.Handle("/admin/editor.wasm.gz", fs)
	http.Handle("/admin/js/", fs)
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
	http.Redirect(w, r, "/file/data/styles.css", 301)
}
