package internal

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("Let's Golang"))

type LoginUser struct {
	Email string
	Token string
}

func init() {
	gob.Register(&LoginUser{})
}

const sessionName = "shizuokago-blog"

func getSessionOptions() *sessions.Options {
	return &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
}

func NewLoginUser(email string, token string) *LoginUser {
	user := LoginUser{}
	user.Email = email
	user.Token = token
	return &user
}

func GetSession(r *http.Request) (*LoginUser, error) {

	sess, err := store.Get(r, sessionName)
	if err != nil {
		return nil, err
	}

	obj := sess.Values["User"]
	if user, ok := obj.(*LoginUser); ok {
		if user.Email == "" {
			return nil, fmt.Errorf("クリアデータ")
		}
		return user, nil
	}
	return nil, fmt.Errorf("ユーザの取得失敗")
}

func SetSession(w http.ResponseWriter, r *http.Request, u *LoginUser) error {

	sess, err := store.Get(r, sessionName)
	if err != nil {
		return err
	}

	if u == nil {
		wk := LoginUser{}
		u = &wk
	}

	sess.Options = getSessionOptions()
	sess.Values["User"] = u

	return sess.Save(r, w)
}
