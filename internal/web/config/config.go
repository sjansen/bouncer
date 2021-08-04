package config

import (
	"context"

	"github.com/sjansen/bouncer/internal/aws"
	"github.com/sjansen/bouncer/internal/config"
	"github.com/sjansen/bouncer/internal/keyring"
)

// Config contains application settings.
type Config struct {
	aws.AWS
	*keyring.KeyRing

	AppURL *config.URL `env:"APP_URL,required"`
	Listen string      `env:"LISTEN"`

	SAML         SAML
	SessionStore SessionStore
}

// SAML contains settings for SAML-based authentication.
type SAML struct {
	EntityID    string `env:"SAML_ENTITY_ID,default=bouncer"`
	MetadataURL string `env:"SAML_METADATA_URL,required"`
	Certificate string `env:"SAML_CERTIFICATE,required"`
	PrivateKey  string `env:"SAML_PRIVATE_KEY,required"`
}

// SessionStore contains setting for app sessions.
type SessionStore struct {
	CreateTable bool       `env:"SESSION_CREATE_TABLE,default=false"`
	EndpointURL config.URL `env:"SESSION_ENDPOINT_URL"`
	TableName   string     `env:"SESSION_TABLE_NAME,required"`
}

func Load(ctx context.Context) (*Config, error) {
	client, err := config.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		KeyRing: keyring.New(client),
	}

	err = config.Load(ctx, cfg)
	if err != nil {
		return nil, err
	}

	aws, err := aws.New(ctx)
	if err != nil {
		return nil, err
	}
	cfg.AWS.Config = aws.Config

	return cfg, nil
}
