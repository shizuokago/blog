package handler

import (
	"log"
	"mime"
	"net/http"

	_ "golang.org/x/tools/playground"
	"golang.org/x/tools/present"
	"golang.org/x/xerrors"

	"github.com/shizuokago/blog/config"
	"github.com/shizuokago/blog/handler/editor"
)

func Register() error {

	present.PlayEnabled = true

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")

	//r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	fs := http.FileServer(http.Dir("./cmd/assets"))
	http.Handle("/js/", fs)
	http.Handle("/css/", fs)
	http.Handle("/images/", fs)
	http.Handle("/favicon.ico", fs)

	//既存のブログが見ている
	http.Handle("/static/css/", fs)

	//Deprecated
	//static := http.FileServer(http.Dir("./cmd"))
	//http.Handle("/static", static)

	// login
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/session", sessionHandler)

	// public page
	http.HandleFunc("/file/", fileHandler)
	http.HandleFunc("/entry/", entryHandler)
	http.HandleFunc("/", topHandler)

	err := editor.Register()
	if err != nil {
		return xerrors.Errorf("editor register: %w", err)
	}
	return nil
}

func Listen() error {

	conf := config.Get()
	s := ":" + conf.Port

	log.Println("Blog Server Start[" + s + "]")

	return http.ListenAndServe(s, nil)
}
