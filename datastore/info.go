package datastore

import (
	"net/http"
)

func PutInformation(r *http.Request) (*User, error) {

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
		return nil, err
	}

	//function
	rtn.Key = getUserKey(r)

	client, err := createClient(c)

	_, err = client.Put(c, rtn.Key, &rtn)
	if err != nil {
		return nil, err
	}
	return &rtn, nil
}
