package editor

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	. "github.com/shizuokago/blog/handler/internal"
)

func Register() error {

	n := mux.NewRouter()
	h := NewHandler(n)

	r := n.PathPrefix("/admin").Subrouter()

	r.HandleFunc("/profile/upload", uploadAvatarHandler).Methods("POST")
	r.HandleFunc("/profile", profileHandler)
	r.HandleFunc("/", adminHandler).Methods("GET")

	r.HandleFunc("/article/create", createArticleHandler).Methods("POST")
	r.HandleFunc("/article/edit/{key}", editArticleHandler).Methods("GET")

	r.HandleFunc("/article/save/{key}", saveArticleHandler).Methods("POST")
	r.HandleFunc("/article/publish/{key}", publishArticleHandler).Methods("POST")
	r.HandleFunc("/article/private/{key}", privateArticleHandler)

	r.HandleFunc("/article/delete/{key}", deleteArticleHandler).Methods("GET")

	r.HandleFunc("/article/bg/save/{key}", saveBackgroundHandler)
	r.HandleFunc("/article/bg/delete/{key}", deleteBackgroundHandler)

	r.HandleFunc("/file/view", viewFileHandler).Methods("GET")
	r.HandleFunc("/file/upload", uploadFileHandler).Methods("POST")
	r.HandleFunc("/file/delete", deleteFileHandler).Methods("POST")
	r.HandleFunc("/file/exists", existsFileHandler).Methods("POST")

	http.Handle("/admin/", h)

	return nil
}

type Handler struct {
	r *mux.Router
}

func NewHandler(r *mux.Router) Handler {
	return Handler{r}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Println("ServeHTTP:" + r.URL.String())
	//セッションの存在を確認
	u, err := GetSession(r)
	if err != nil {
		http.Redirect(w, r, "/login", 301)
		return
	}

	if u == nil {
		log.Println("ユーザがいない")
		http.Redirect(w, r, "/login", 301)
		return
	}

	h.r.ServeHTTP(w, r)
}

func deleteDir(s string) string {
	ds := []byte(s)
	return string(ds[5:])
}
