package internal

import (
	"embed"
	"io/fs"
	"log"
	"mime"
	"net/http"
)

//go:embed _assets/static
var embStatic embed.FS
var staticFS fs.FS

func init() {
	var err error
	staticFS, err = fs.Sub(embStatic, "_assets/static")
	if err != nil {
		log.Printf("fs.Sub() error: %+v", err)
	}
}

func RegisterStatic() error {

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")

	fs := http.FileServer(http.FS(staticFS))

	http.Handle("/favicon.ico", fs)
	http.HandleFunc("/admin/bin/editor.wasm.gz", WasmServer(http.FS(staticFS)))
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
