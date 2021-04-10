package internal

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/rakyll/statik/fs"
	"golang.org/x/xerrors"
)

var templateMap map[string]string

func init() {
}

func GetTemplate(funcMap template.FuncMap, path ...string) (*template.Template, error) {

	var err error
	var rtn *template.Template
	for _, elm := range path {
		buf, ok := templateMap["/"+elm]
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
	err := fs.Walk(statikFS, "/templates/", func(path string, fi os.FileInfo, err error) error {

		r, err := statikFS.Open(path)
		if err != nil {
			return xerrors.Errorf("statik Open() error: %w", err)
		}
		defer r.Close()

		byt, err := ioutil.ReadAll(r)
		if err != nil {
			return xerrors.Errorf("ioutil.ReadAll() error: %w", err)
		}

		templateMap[strings.Replace(path, "/templates", "", 1)] = string(byt)
		return nil
	})

	if err != nil {
		return xerrors.Errorf("fs.Walk() error: %w", err)
	}
	return nil
}
