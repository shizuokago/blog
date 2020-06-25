package datastore

import (
	"net/http"
)

type FileDs struct {
	request *http.Request
}

func (ds FileDs) readFile(name string) ([]byte, error) {
	key := "data/" + name
	file, err := GetFileData(ds.request, key)
	if err != nil {
		return []byte(err.Error()), err
	}
	return file.Content, nil
}
