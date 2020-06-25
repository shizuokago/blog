package datastore

import (
	"errors"
	"net/http"

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"
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

const KIND_FILE = "File"

type File struct {
	Size int64
	Type int64
	Meta
}

const (
	FILE_TYPE_BG     = 1
	FILE_TYPE_AVATAR = 2
	FILE_TYPE_DATA   = 3
)

func getFileKey(r *http.Request, name string) *datastore.Key {
	return datastore.NameKey(KIND_FILE, name, nil)
}

func SelectFile(r *http.Request, p int) ([]File, error) {

	c := r.Context()

	q := datastore.NewQuery(KIND_FILE).
		Filter("Type =", 3).
		Order("- UpdatedAt").
		Limit(10)

	var s []File

	client, err := createClient(c)
	if err != nil {
		return nil, err
	}
	t := client.Run(c, q)
	for {
		var f File
		key, err := t.Next(&f)

		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		f.SetKey(key)
		s = append(s, f)
	}

	return s, nil
}

func DeleteFile(r *http.Request, id string) error {

	c := r.Context()

	fkey := getFileKey(r, id)

	client, err := createClient(c)
	err = client.Delete(c, fkey)
	if err != nil {
		return err
	}

	fdkey := getFileDataKey(r, id)
	err = client.Delete(c, fdkey)
	return err
}

func ExistsFile(r *http.Request, id string, t int64) (bool, error) {

	c := r.Context()
	dir := "data"
	if t == FILE_TYPE_BG {
		dir = "bg"
	} else if t == FILE_TYPE_AVATAR {
		dir = "avatar"
	}

	key := getFileKey(r, dir+"/"+id)

	rtn := File{}
	client, err := createClient(c)
	err = client.Get(c, key, &rtn)
	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return false, nil
		} else {
			return true, err
		}
	}

	return true, nil

}

func SaveFile(r *http.Request, id string, t int64) error {

	upload, header, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer upload.Close()

	b, flg, err := ConvertImage(upload)
	if err != nil {
		return err
	}

	c := r.Context()

	dir := "data"
	if t == FILE_TYPE_BG {
		dir = "bg"
	} else if t == FILE_TYPE_AVATAR {
		dir = "avatar"
	} else {
		if id == "" {
			id = header.Filename
		}
	}

	fid := dir + "/" + id
	file := &File{
		Size: int64(len(b)),
		Type: t,
	}

	file.Key = getFileKey(r, fid)

	client, err := createClient(c)
	_, err = client.Put(c, file.Key, file)
	if err != nil {
		return err
	}

	mime := header.Header["Content-Type"][0]
	if flg {
		mime = "image/jpeg"
	}

	fileData := &FileData{
		Content: b,
		Mime:    mime,
	}
	fileData.SetKey(getFileDataKey(r, fid))
	_, err = client.Put(c, fileData.GetKey(), fileData)
	if err != nil {
		return err
	}
	return nil
}

func SaveBackgroundImage(r *http.Request, id string) error {
	err := SaveFile(r, id, FILE_TYPE_BG)
	return err
}

func DeleteBackgroundImage(r *http.Request, id string) error {
	err := DeleteFile(r, "bg/"+id)
	return err
}

const KIND_FILEDATA = "FileData"

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

func getFileDataKey(r *http.Request, name string) *datastore.Key {
	return datastore.NameKey(KIND_FILEDATA, name, nil)
}

func GetFileData(r *http.Request, name string) (*FileData, error) {
	c := r.Context()

	rtn := FileData{}
	key := getFileDataKey(r, name)

	client, err := createClient(c)
	err = client.Get(c, key, &rtn)
	if err != nil {
		if errors.Is(err, datastore.ErrNoSuchEntity) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &rtn, nil
}
