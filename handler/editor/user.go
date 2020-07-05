package editor

import (
	"net/http"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {

	user, err := GetSession(r)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}

	var u *datastore.User
	if r.Method == "POST" {
		u, err = datastore.PutInformation(r, user.Email)
		if err != nil {
			ErrorPage(w, "InternalServerError", err, 500)
			return
		}
	}

	u, err = datastore.GetUser(r, user.Email)
	if u == nil {
		u = &datastore.User{}
	}

	bgd := datastore.GetBlog(r)
	data := struct {
		Blog *datastore.Blog
		User *datastore.User
	}{bgd, u}

	adminRender(w, "./cmd/templates/admin/profile.tmpl", data)
}

func uploadAvatarHandler(w http.ResponseWriter, r *http.Request) {

	u, err := GetSession(r)

	err = datastore.SaveAvatar(r, u.Email)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/profile", 301)
}
