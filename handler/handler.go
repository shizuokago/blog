package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"mime"
	"net/http"

	_ "golang.org/x/tools/playground"
	"golang.org/x/tools/present"
	"golang.org/x/xerrors"

	"github.com/shizuokago/blog/config"
	"github.com/shizuokago/blog/handler/editor"
	. "github.com/shizuokago/blog/handler/internal"
)

func Register() error {

	present.PlayEnabled = true

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")

	//r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	fs := http.FileServer(http.Dir("./cmd/static"))
	http.Handle("/js/", fs)
	http.Handle("/css/", fs)
	http.Handle("/static/css/", fs) // 既存ブログ用に設定
	http.Handle("/images/", fs)
	http.Handle("/favicon.ico", fs)

	http.HandleFunc("/entry/", entryHandler)
	http.HandleFunc("/file/", fileHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/session", sessionHandler)

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

	log.Println("Start Blog Server[" + s + "]")

	return http.ListenAndServe(s, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	err := SetSession(w, r, nil)
	if err != nil {
	}

	tmpl, err := template.ParseFiles("cmd/static/templates/authentication.tmpl")
	if err != nil {
		log.Printf("Error Page Parse Error: %v", err)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error Page Execute Error: %v", err)
		return
	}

}

func sessionHandler(w http.ResponseWriter, r *http.Request) {

	code := 200
	dto := struct {
		Success bool
	}{false}

	r.ParseForm()
	email := r.FormValue("email")
	token := r.FormValue("token")

	//TODO
	log.Println(email)
	flag := true

	dto.Success = flag

	if !flag {
		//403を返す
		code = 403
	} else {
		//Cookieの作成
		u := NewLoginUser(email, token)

		err := SetSession(w, r, u)
		if err != nil {
			code = 500
			dto.Success = false
			log.Println(err)
		}
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(dto)

}
