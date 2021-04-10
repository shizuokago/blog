package internal

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"

	"golang.org/x/xerrors"
)

//go:embed _assets/templates
var embTmpls embed.FS
var tmpls fs.FS

var templateMap map[string]string

func init() {

	var err error
	tmpls, err = fs.Sub(embTmpls, "_assets/templates")
	if err != nil {
		log.Printf("fs.Sub() error: %+v", err)
		return
	}

	err = initTemplates()
	if err != nil {
		log.Printf("initTemplate() error: %+v", err)
	}
}

func GetTemplate(funcMap template.FuncMap, path ...string) (*template.Template, error) {

	var err error
	var rtn *template.Template
	for _, elm := range path {
		buf, ok := templateMap[elm]
		if !ok {
			return nil, fmt.Errorf("unknown template[%s]", elm)
		}
		if rtn == nil {
			rtn = template.New("root")
			if funcMap != nil {
				rtn = rtn.Funcs(funcMap)
			}
		}
		rtn, err = rtn.Parse(buf)
		if err != nil {
			return nil, xerrors.Errorf("Parse template[%s]]: %w", elm, err)
		}
	}
	return rtn, nil
}

func initTemplates() error {

	templateMap = make(map[string]string)

	err := setTemplate(".")
	if err != nil {
		return xerrors.Errorf("error: %w", err)
	}
	return nil
}

func setTemplate(path string) error {

	dirs, err := fs.ReadDir(tmpls, path)
	for _, dir := range dirs {

		name := dir.Name()
		read := name

		if path != "." {
			read = path + "/" + name
		}

		if dir.IsDir() {
			err = setTemplate(read)
			if err != nil {
				return xerrors.Errorf("error: %w", err)
			}

		} else {

			r, err := tmpls.Open(read)
			if err != nil {
				return xerrors.Errorf("filesystem open[%s] : %w", path, err)
			}

			h, err := io.ReadAll(r)
			if err != nil {
				return xerrors.Errorf("data read[%s] : %w", path, err)
			}
			templateMap[read] = string(h)
		}
	}

	return nil

}
