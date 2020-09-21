package editor

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/shizuokago/blog/datastore"
	. "github.com/shizuokago/blog/handler/internal"
)

func viewCacheHandler(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	caches, err := datastore.GetCaches(ctx)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	dto := struct {
		Caches []*datastore.Cache
	}{caches}

	adminRender(w, "./cmd/templates/admin/cache.tmpl", dto)
}

func registerCacheHandler(w http.ResponseWriter, r *http.Request) {

	p := r.FormValue("pattern")
	if p == "" || strings.Index(p, "/") == -1 {
		ErrorPage(w, "Internal Server Error", fmt.Errorf("URL pattern error"), 500)
		return
	}

	ageBuf := r.FormValue("age")
	age, err := strconv.Atoi(ageBuf)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	ctx := r.Context()

	_, err = datastore.RegisterCache(ctx, p, age)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	http.Redirect(w, r, "/admin/cache/view", 302)
}

func updateCacheHandler(w http.ResponseWriter, r *http.Request) {

	keyBuf := r.FormValue("key")

	p := r.FormValue("pattern")
	if p == "" || strings.Index(p, "/") == -1 {
		ErrorPage(w, "Internal Server Error", fmt.Errorf("URL pattern error"), 500)
		return
	}

	ageBuf := r.FormValue("age")
	age, err := strconv.Atoi(ageBuf)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	var cache datastore.Cache
	key := datastore.GetCacheKey(keyBuf)

	ctx := r.Context()

	err = datastore.Get(ctx, key, &cache)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	cache.Pattern = p
	cache.MaxAge = age

	err = datastore.Put(ctx, &cache)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	http.Redirect(w, r, "/admin/cache/view", 302)
}

func deleteCacheHandler(w http.ResponseWriter, r *http.Request) {

	keyBuf := r.FormValue("key")
	ctx := r.Context()
	key := datastore.GetCacheKey(keyBuf)

	err := datastore.Delete(ctx, key)
	if err != nil {
		ErrorPage(w, "Internal Server Error", err, 500)
		return
	}

	http.Redirect(w, r, "/admin/cache/view", 302)
}
