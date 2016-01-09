package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

func fileHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["key"]

	file, err := getFileData(r, name)
	if err != nil {
	}

	//set MIME

	w.Write(file.Content)
}
