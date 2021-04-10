package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	. "github.com/shizuokago/blog/handler/internal"
	"golang.org/x/xerrors"
)

func registerLogin() error {
	// login
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/session", sessionHandler)
	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	err := SetSession(w, r, nil)
	if err != nil {
		log.Printf("Error set session Error: %v", err)
	}

	err = render(w, nil, "authentication.tmpl")
	if err != nil {
		log.Printf("render() Error: %v", err)
	}

}

func render(w http.ResponseWriter, dto interface{}, name string) error {
	funcMap := template.FuncMap{}

	tmpl, err := GetTemplate(funcMap, name)
	if err != nil {
		return xerrors.Errorf("GetTemplate() error: %w", err)
	}
	err = tmpl.Execute(w, dto)
	if err != nil {
		return xerrors.Errorf("template.Execute() error: %w", err)
	}
	return nil
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := SetSession(w, r, nil)
	if err != nil {
		log.Printf("Error set session Error: %v", err)
	}
	http.Redirect(w, r, "/", 302)
}

func sessionHandler(w http.ResponseWriter, r *http.Request) {

	code := 200
	dto := struct {
		Success bool
	}{false}

	r.ParseForm()
	email := r.FormValue("email")
	token := r.FormValue("token")

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

	encoder := json.NewEncoder(w)
	err := encoder.Encode(dto)
	if err != nil {
		log.Println(err)
	}
}
