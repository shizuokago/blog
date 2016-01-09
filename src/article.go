package main

import (
	"net/http"
)

func createArticleHandler(w http.ResponseWriter, r *http.Request) {

	createArticle(r)

	// Render Editor
}
