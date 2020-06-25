package handler

import (
	"encoding/json"
	"mime"
	"net/http"

	"github.com/gorilla/mux"
	_ "golang.org/x/tools/playground"
	"golang.org/x/tools/present"

	"github.com/shizuokago/blog/config"
)

func Register() error {

	present.PlayEnabled = true

	// App Engine has no /etc/mime.types
	mime.AddExtensionType(".svg", "image/svg+xml")

	r := mux.NewRouter()
	r.HandleFunc("/admin/profile/upload", uploadAvatarHandler).Methods("POST")
	r.HandleFunc("/admin/profile", profileHandler)
	r.HandleFunc("/admin/", adminHandler).Methods("GET")

	r.HandleFunc("/admin/article/create", createArticleHandler).Methods("POST")
	r.HandleFunc("/admin/article/edit/{key}", editArticleHandler).Methods("GET")

	r.HandleFunc("/admin/article/save/{key}", saveArticleHandler).Methods("POST")
	r.HandleFunc("/admin/article/publish/{key}", publishArticleHandler).Methods("POST")
	r.HandleFunc("/admin/article/private/{key}", privateArticleHandler)

	r.HandleFunc("/admin/article/delete/{key}", deleteArticleHandler).Methods("GET")

	r.HandleFunc("/admin/article/bg/save/{key}", saveBackgroundHandler)
	r.HandleFunc("/admin/article/bg/delete/{key}", deleteBackgroundHandler)

	r.HandleFunc("/admin/file/view", viewFileHandler).Methods("GET")
	r.HandleFunc("/admin/file/upload", uploadFileHandler).Methods("POST")
	r.HandleFunc("/admin/file/delete", deleteFileHandler).Methods("POST")
	r.HandleFunc("/admin/file/exists", existsFileHandler).Methods("POST")

	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	r.HandleFunc("/", topHandler).Methods("GET")
	r.HandleFunc("/entry/{key}", entryHandler).Methods("GET")
	http.HandleFunc("/file/", fileHandler)
	http.Handle("/", r)

	return nil
}

func Listen() error {

	conf := config.Get()

	return http.ListenAndServe(":"+conf.Port, nil)
}

func deleteDir(s string) string {
	ds := []byte(s)
	return string(ds[5:])
}

func errorPage(w http.ResponseWriter, t, m string, code int) {

	data := struct {
		Code    int
		Title   string
		Message string
	}{code, t, m}

	w.WriteHeader(data.Code)
	err := errorTmpl.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

func errorJson(w http.ResponseWriter, t, m string, code int) {

	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	d := map[string]interface{}{
		"success": false,
		"title":   t,
		"message": m,
	}
	enc.Encode(d)
}
