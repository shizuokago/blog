package blog

import (
	"golang.org/x/xerrors"

	"github.com/shizuokago/blog/config"
	"github.com/shizuokago/blog/handler"
)

func Start(opts ...config.Option) error {

	err := config.Set(opts...)
	if err != nil {
		return xerrors.Errorf("config set: %w", err)
	}

	err = handler.Register()
	if err != nil {
		return xerrors.Errorf("handler register: %w", err)
	}

	err = handler.Listen()
	if err != nil {
		return xerrors.Errorf("handler listen: %w", err)
	}

	return nil
}
