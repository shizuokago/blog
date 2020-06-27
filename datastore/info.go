package datastore

import (
	"net/http"

	"golang.org/x/xerrors"
)

func PutInformation(r *http.Request, key string) (*User, error) {

	c := r.Context()

	r.ParseForm()

	save := false
	if r.FormValue("AutoSave") != "" {
		save = true
	}

	rtn := User{
		Name:      r.FormValue("Name"),
		Job:       r.FormValue("Job"),
		Email:     r.FormValue("Email"),
		URL:       r.FormValue("Url"),
		TwitterId: r.FormValue("TwitterId"),
		AutoSave:  save,
	}

	err := PutBlog(r)
	if err != nil {
		return nil, xerrors.Errorf("put blog: %w", err)
	}

	//function
	rtn.Key = getUserKey(key)

	err = Put(c, &rtn)
	if err != nil {
		return nil, xerrors.Errorf("put user: %w", err)
	}
	return &rtn, nil
}
