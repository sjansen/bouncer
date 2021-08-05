package config

import (
	"context"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/sjansen/bouncer/internal/aws"
)

type ssmClient interface {
	DescribeParameters(ctx context.Context, names ...string) (map[string]time.Time, error)
	GetParameter(ctx context.Context, name string) (string, error)
	GetParameters(ctx context.Context, names ...string) (map[string]string, map[string]time.Time, error)
	PutParameter(ctx context.Context, key, value string, encrypt bool) error
}

type ssmConfig struct {
	Prefix string `env:"SSM_PREFIX,required"`
}

type Client struct {
	ssmClient
}

func NewClient(ctx context.Context) (*Client, error) {
	cfg := &ssmConfig{}
	if err := load(ctx, cfg, envconfig.OsLookuper()); err != nil {
		return nil, err
	}

	aws, err := aws.New(ctx)
	if err != nil {
		return nil, err
	}

	client, err := aws.NewSSMClient(cfg.Prefix)
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

func (c *Client) DescribeJWKSet(ctx context.Context) (map[string]time.Time, error) {
	return c.ssmClient.DescribeParameters(ctx, "JWK1", "JWK2")
}

func (c *Client) GetJWKSet(ctx context.Context) (map[string]string, map[string]time.Time, error) {
	return c.ssmClient.GetParameters(ctx, "JWK1", "JWK2")
}

func (c *Client) PutJWK(ctx context.Context, name, value string) error {
	return c.ssmClient.PutParameter(ctx, name, value, true)
}

func (c *Client) PutSAMLCertificate(ctx context.Context, value string) error {
	return c.ssmClient.PutParameter(ctx, "SAML_CERTIFICATE", value, false)
}

func (c *Client) PutSAMLPrivateKey(ctx context.Context, value string) error {
	return c.ssmClient.PutParameter(ctx, "SAML_PRIVATE_KEY", value, true)
}
