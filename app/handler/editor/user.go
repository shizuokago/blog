package editor

import (
	"net/http"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
	. "github.com/shizuokago/blog/handler/internal/form"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {

	lg, err := GetSession(r)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}

	var u *datastore.User
	if r.Method == "POST" {

		ctx := r.Context()

		blog, err := CreateBlog(r)
		if err != nil {
			ErrorPage(w, "InternalServerError", err, 500)
			return
		}

		user, err := CreateUser(r)
		if err != nil {
			ErrorPage(w, "InternalServerError", err, 500)
			return
		}

		err = datastore.PutInformation(ctx, blog, user, lg.Email)
		if err != nil {
			ErrorPage(w, "InternalServerError", err, 500)
			return
		}
	}

	ctx := r.Context()
	u, err = datastore.GetUser(ctx, lg.Email)
	if u == nil {
		u = &datastore.User{}
	}

	bgd := datastore.GetBlog(ctx)
	data := struct {
		Blog *datastore.Blog
		User *datastore.User
	}{bgd, u}

	adminRender(w, "profile.tmpl", data)
}

func uploadAvatarHandler(w http.ResponseWriter, r *http.Request) {

	u, err := GetSession(r)

	ctx := r.Context()

	f, d, err := GetFile(r)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	p := datastore.FileParam{
		File:     f,
		FileData: d,
	}

	err = datastore.SaveAvatar(ctx, u.Email, &p)
	if err != nil {
		ErrorPage(w, "InternalServerError", err, 500)
		return
	}
	http.Redirect(w, r, "/admin/profile", 302)
}
