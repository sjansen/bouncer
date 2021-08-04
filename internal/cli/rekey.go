package cli

import (
	"context"

	"github.com/sjansen/bouncer/internal/config"
	"github.com/sjansen/bouncer/internal/keyring"
)

type rekeyCmd struct {
	JWK  rekeyJWK  `kong:"cmd,name='jwk'"`
	SAML rekeySAML `kong:"cmd,name='saml'"`
}

type rekeyJWK struct{}

func (cmd *rekeyJWK) Run() error {
	ctx := context.Background()
	client, err := config.NewClient(ctx)
	if err != nil {
		return err
	}

	keyring := keyring.New(client)
	return keyring.RotateJWKs(ctx)
}

type rekeySAML struct{}

func (cmd *rekeySAML) Run() error {
	ctx := context.Background()
	client, err := config.NewClient(ctx)
	if err != nil {
		return err
	}

	keyring := keyring.New(client)
	return keyring.RekeySAML(ctx)
}
