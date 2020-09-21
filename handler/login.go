package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	. "github.com/shizuokago/blog/handler/internal"
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

	tmpl, err := template.ParseFiles("cmd/templates/authentication.tmpl")
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

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	err := SetSession(w, r, nil)
	if err != nil {
		log.Printf("Error set session Error: %v", err)
	}
	http.Redirect(w, r, "/", 301)
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
