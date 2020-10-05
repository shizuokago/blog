package form

import (
	"net/http"

	"github.com/shizuokago/blog/datastore"
	"github.com/shizuokago/blog/logic"
	"golang.org/x/xerrors"
)

func GetFile(r *http.Request) (*datastore.File, *datastore.FileData, error) {

	upload, header, err := r.FormFile("file")
	if err != nil {
		return nil, nil, xerrors.Errorf("from file: %w", err)
	}
	defer upload.Close()

	b, flg, err := logic.ConvertImage(upload)
	if err != nil {
		return nil, nil, xerrors.Errorf("convert image: %w", err)
	}

	file := &datastore.File{
		Size: int64(len(b)),
	}

	mime := header.Header["Content-Type"][0]
	if flg {
		mime = "image/jpeg"
	}

	fileData := &datastore.FileData{
		Content: b,
		Mime:    mime,
	}

	return file, fileData, nil
}
