package form

import (
	"net/http"
	"strings"

	"github.com/shizuokago/blog/datastore"
)

func CreateArticle(r *http.Request) (*datastore.Article, error) {

	r.ParseForm()

	title := r.FormValue("Title")
	tags := r.FormValue("Tags")
	mark := r.FormValue("Markdown")

	art := datastore.Article{}

	art.Title = title
	art.SubTitle = CreateSubTitle(mark)
	art.Tags = tags
	art.Markdown = []byte(mark)

	return &art, nil
}

func CreateNewArticle(b *datastore.Blog) *datastore.Article {

	base := b.Template
	article := &datastore.Article{
		Title:    "New Title",
		Tags:     b.Tags,
		Markdown: []byte(base),
	}
	return article
}

func CreateSubTitle(src string) string {

	dst := strings.Replace(src, "\n", "", -1)
	dst = strings.Replace(dst, "*", "", -1)

	if len(dst) > 600 {
		dst = string([]rune(dst)[0:200]) + "..."
	}
	return dst
}
