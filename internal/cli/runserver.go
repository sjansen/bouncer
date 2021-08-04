package cli

import (
	"context"

	"github.com/sjansen/bouncer/internal/web/server"
)

type runserverCmd struct{}

func (cmd *runserverCmd) Run() error {
	ctx := context.Background()
	s, err := server.New(ctx)
	if err != nil {
		return err
	}

	return s.ListenAndServe()
}
