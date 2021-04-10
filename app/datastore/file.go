package datastore

import (
	"context"
	"errors"
	"log"

	"cloud.google.com/go/datastore"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
)

var fileCursor map[int]string

type FileParam struct {
	File     *File
	FileData *FileData
}

func init() {
	fileCursor = make(map[int]string)
}

const KindFile = "File"

type File struct {
	Size int64
	Type int64
	Meta
}

const (
	FileTypeBG     = 1
	FileTypeAvatar = 2
	FileTypeData   = 3
)

func getFileKey(name string) *datastore.Key {
	return datastore.NameKey(KindFile, name, nil)
}

func SelectFile(ctx context.Context, p int) ([]File, error) {

	q := datastore.NewQuery(KindFile).
		Filter("Type =", 3).
		Order("- UpdatedAt").
		Limit(10)

	var s []File

	client, err := createClient(ctx)
	if err != nil {
		return nil, xerrors.Errorf("create client: %w", err)
	}

	if cur, ok := fileCursor[p]; ok {
		cursor, err := datastore.DecodeCursor(cur)
		if err != nil {
			log.Printf("%+v", err)
		} else {
			q = q.Start(cursor)
		}
	}

	t := client.Run(ctx, q)
	for {
		var f File
		key, err := t.Next(&f)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf("next: %w", err)
		}
		f.SetKey(key)
		s = append(s, f)
	}

	next, err := t.Cursor()
	if err != nil {
		log.Printf("%+v", err)
	} else {
		fileCursor[p+1] = next.String()
	}
	return s, nil
}

func DeleteFile(ctx context.Context, id string) error {

	fkey := getFileKey(id)
	client, err := createClient(ctx)
	if err != nil {
		return xerrors.Errorf("create client: %w", err)
	}

	err = client.Delete(ctx, fkey)
	if err != nil {
		return xerrors.Errorf("delete file: %w", err)
	}

	fdkey := getFileDataKey(id)
	err = client.Delete(ctx, fdkey)
	if err != nil {
		return xerrors.Errorf("delete file data: %w", err)
	}

	fileCursor = make(map[int]string)
	return nil
}

func ExistsFile(ctx context.Context, id string, t int64) (bool, error) {

	dir := "data"
	if t == FileTypeBG {
		dir = "bg"
	} else if t == FileTypeAvatar {
		dir = "avatar"
	}

	key := getFileKey(dir + "/" + id)

	rtn := File{}
	err := Get(ctx, key, &rtn)

	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return false, nil
		} else {
			return true, xerrors.Errorf("file get: %w", err)
		}
	}

	return true, nil

}

func SaveFile(ctx context.Context, id string, t int64, file *FileParam) error {

	dir := "data"
	if t == FileTypeBG {
		dir = "bg"
	} else if t == FileTypeAvatar {
		dir = "avatar"
	}

	if id == "" {
		return errors.New("id is empty")
	}

	fid := dir + "/" + id
	file.File.SetKey(getFileKey(fid))
	file.File.Type = t

	file.FileData.SetKey(getFileDataKey(fid))

	err := Put(ctx, file.File)
	if err != nil {
		return xerrors.Errorf("put file: %w", err)
	}

	err = Put(ctx, file.FileData)
	if err != nil {
		return xerrors.Errorf("put file data: %w", err)
	}

	fileCursor = make(map[int]string)
	return nil
}

func SaveBackgroundImage(ctx context.Context, id string, file *FileParam) error {
	err := SaveFile(ctx, id, FileTypeBG, file)
	if err != nil {
		return xerrors.Errorf("save file: %w", err)
	}
	return nil
}

func DeleteBackgroundImage(ctx context.Context, id string) error {
	err := DeleteFile(ctx, "bg/"+id)
	if err != nil {
		return xerrors.Errorf("delete background file: %w", err)
	}
	return nil
}

const KindFileData = "FileData"

type FileData struct {
	key     *datastore.Key
	Mime    string
	Content []byte `datastore:",noindex"`
}

func (d *FileData) GetKey() *datastore.Key {
	return d.key
}

func (d *FileData) SetKey(k *datastore.Key) {
	d.key = k
}

func getFileDataKey(name string) *datastore.Key {
	return datastore.NameKey(KindFileData, name, nil)
}

func GetFileData(ctx context.Context, name string) (*FileData, error) {

	rtn := FileData{}
	key := getFileDataKey(name)

	err := Get(ctx, key, &rtn)
	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		} else {
			return nil, xerrors.Errorf("get file data: %w", err)
		}
	}

	return &rtn, nil
}
