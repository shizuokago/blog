package handler

import (
	"log"
	"net/http"

	_ "golang.org/x/tools/playground"
	"golang.org/x/tools/present"
	"golang.org/x/xerrors"

	"github.com/shizuokago/blog/config"
	"github.com/shizuokago/blog/handler/editor"
)

func Register() error {

	present.PlayEnabled = true

	err := registerStatic()
	if err != nil {
		return xerrors.Errorf("static register: %w", err)
	}

	err = editor.Register()
	if err != nil {
		return xerrors.Errorf("editor register: %w", err)
	}

	err = registerLogin()
	if err != nil {
		return xerrors.Errorf("login register: %w", err)
	}

	err = registerPublic()
	if err != nil {
		return xerrors.Errorf("public register: %w", err)
	}

	return nil
}

func Listen() error {

	conf := config.Get()
	s := ":" + conf.Port

	log.Println("Blog Server Start[" + s + "]")

	return http.ListenAndServe(s, nil)
}
