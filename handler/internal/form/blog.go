package form

import (
	"net/http"

	"github.com/shizuokago/blog/datastore"
)

func CreateBlog(r *http.Request) (*datastore.Blog, error) {

	blog := datastore.Blog{
		Name:        r.FormValue("BlogName"),
		Author:      r.FormValue("BlogAuthor"),
		Description: r.FormValue("Description"),
		Tags:        r.FormValue("BlogTags"),
		Template:    r.FormValue("BlogTemplate"),
		Users:       r.FormValue("Users"),
	}

	return &blog, nil
}
