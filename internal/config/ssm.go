package config

import (
	"context"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
	"github.com/sjansen/bouncer/internal/aws"
)

type ssmClient interface {
	DescribeParameters(ctx context.Context, names ...string) (map[string]time.Time, error)
	GetParameter(ctx context.Context, name string) (string, error)
	PutParameter(ctx context.Context, key, value string, encrypt bool) error
}

type ssmConfig struct {
	Prefix string `env:"SSM_PREFIX,required"`
}

type ssmMutator struct {
	ssmClient
}

func newSSMClient(ctx context.Context) (*aws.SSMClient, error) {
	cfg := &ssmConfig{}
	if err := load(ctx, cfg, envconfig.OsLookuper()); err != nil {
		return nil, err
	}

	aws, err := aws.New(ctx)
	if err != nil {
		return nil, err
	}

	return aws.NewSSMClient(cfg.Prefix)
}

func (m *ssmMutator) resolve(ctx context.Context, k, v string) (string, error) {
	v = strings.TrimSpace(v)
	if v == "ssm" {
		return m.ssmClient.GetParameter(ctx, k)
	}
	return v, nil
}
