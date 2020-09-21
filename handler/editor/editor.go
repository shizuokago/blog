package editor

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

func Register() error {

	n := mux.NewRouter()

	r := n.PathPrefix("/admin").Subrouter()
	h := NewLoginHandler(r)

	r.HandleFunc("/profile/upload", uploadAvatarHandler).Methods("POST")
	r.HandleFunc("/profile", profileHandler)

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

	r.HandleFunc("/cache/view", viewCacheHandler).Methods("GET")
	r.HandleFunc("/cache/register", registerCacheHandler).Methods("POST")
	r.HandleFunc("/cache/update", updateCacheHandler).Methods("POST")
	r.HandleFunc("/cache/delete", deleteCacheHandler).Methods("POST")

	r.HandleFunc("/", adminHandler).Methods("GET")

	http.Handle("/admin/", h)

	return nil
}

type LoginHandler struct {
	r *mux.Router
}

func NewLoginHandler(r *mux.Router) LoginHandler {
	return LoginHandler{r}
}

func (h LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	u, err := GetSession(r)
	if err != nil {
		log.Printf("session : %+v", err)
		http.Redirect(w, r, "/login", 302)
		return
	}

	if u == nil {
		log.Println("ユーザがいない")
		http.Redirect(w, r, "/login", 302)
		return
	}

	//メンバ設定
	if !datastore.IsUser(r, u.Email) {
		log.Println("ユーザが違う:" + u.Email)
		http.Redirect(w, r, "/logout", 302)
		return
	}

	h.r.ServeHTTP(w, r)
}
