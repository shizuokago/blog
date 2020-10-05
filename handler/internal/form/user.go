package form

import (
	"net/http"

	"github.com/shizuokago/blog/datastore"
)

func CreateUser(r *http.Request) (*datastore.User, error) {

	r.ParseForm()

	save := false
	if r.FormValue("AutoSave") != "" {
		save = true
	}

	rtn := datastore.User{
		Name:      r.FormValue("Name"),
		Job:       r.FormValue("Job"),
		Email:     r.FormValue("Email"),
		URL:       r.FormValue("Url"),
		TwitterId: r.FormValue("TwitterId"),
		AutoSave:  save,
	}

	return &rtn, nil
}
