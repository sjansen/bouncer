package cli

import (
	"context"

	"github.com/sjansen/bouncer/internal/config"
)

type rekeyCmd struct {
	JWK  rekeyJWK  `kong:"cmd,name='jwk'"`
	SAML rekeySAML `kong:"cmd,name='saml'"`
}

type rekeyJWK struct{}

func (cmd *rekeyJWK) Run() error {
	ctx := context.Background()
	return config.RotateJWKs(ctx)
}

type rekeySAML struct{}

func (cmd *rekeySAML) Run() error {
	ctx := context.Background()
	return config.RekeySAML(ctx)
}
