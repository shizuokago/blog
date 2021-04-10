package logic

import (
	"context"

	"github.com/shizuokago/blog/datastore"
)

type FileDs struct {
	ctx context.Context
}

func (ds FileDs) readFile(name string) ([]byte, error) {
	key := "data/" + name
	file, err := datastore.GetFileData(ds.ctx, key)
	if err != nil {
		return []byte(err.Error()), err
	}
	return file.Content, nil
}
