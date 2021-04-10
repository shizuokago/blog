package internal

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/shizuokago/blog/datastore"
	"golang.org/x/xerrors"
)

var cacheMap map[string]*datastore.Cache

func init() {
	cacheMap = make(map[string]*datastore.Cache)
}

func AddCacheHeader(w http.ResponseWriter, r *http.Request) error {

	if len(cacheMap) == 0 {
		err := reloadCaches(r)
		if err != nil {
			fmt.Printf("reload caches: %+v\n", err)
		}
	}

	cache := searchCache(r)
	if cache != nil {
		h := w.Header()
		h.Set("Cache-Control", fmt.Sprintf("public, max-age=%d", cache.MaxAge))
	}

	return nil
}

func searchCache(r *http.Request) *datastore.Cache {

	url := r.URL.String()

	for k, v := range cacheMap {
		if matchURL(url, k) {
			return v
		}
	}

	return nil
}

func matchURL(src, v string) bool {

	srcS := strings.Split(src, "/")
	vS := strings.Split(v, "/")
	if len(srcS) != len(vS) {
		return false
	}

	for idx, elm := range srcS {
		oppo := vS[idx]
		if oppo == "*" {
			continue
		}
		if elm == oppo {
			continue
		}
		return false
	}

	return true
}

func reloadCaches(r *http.Request) error {

	caches, err := datastore.GetCaches(r.Context())
	if err != nil {
		return xerrors.Errorf("get caches: %w", err)
	}

	sort.Slice(caches, func(i, j int) bool {
		c1 := caches[i]
		c2 := caches[j]

		s1 := strings.Split(c1.Pattern, "/")
		s2 := strings.Split(c2.Pattern, "/")

		if len(s1) > len(s2) {
			return true
		} else if len(s1) > len(s2) {
			return false
		}

		if len(c1.Pattern) > len(c2.Pattern) {
			return true
		} else if len(c1.Pattern) > len(c2.Pattern) {
			return false
		}

		return true
	})

	for _, elm := range caches {
		cacheMap[elm.Pattern] = elm
	}

	return nil
}
