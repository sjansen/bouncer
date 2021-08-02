package cli

import (
	"context"

	"github.com/sjansen/bouncer/internal/config"
)

type rekeyCmd struct{}

func (cmd *rekeyCmd) Run() error {
	ctx := context.Background()
	return config.Rekey(ctx)
}
