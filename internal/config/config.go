package config

import (
	"context"
	"strings"

	"github.com/sethvargo/go-envconfig"
)

// Load reads settings from the environment and AWS SSM.
func Load(ctx context.Context, config interface{}) error {
	client, err := NewClient(ctx)
	if err != nil {
		return err
	}

	m := &ssmMutator{client}
	return load(ctx, config, envconfig.OsLookuper(), m.resolve)
}

func load(ctx context.Context, config interface{}, l envconfig.Lookuper, fns ...envconfig.MutatorFunc) error {
	err := envconfig.ProcessWith(ctx, config,
		envconfig.PrefixLookuper("BOUNCER_", l),
		fns...,
	)
	return err
}

type ssmMutator struct {
	*Client
}

func (m *ssmMutator) resolve(ctx context.Context, k, v string) (string, error) {
	v = strings.TrimSpace(v)
	if v == "ssm" {
		return m.Client.GetParameter(ctx, k)
	}
	return v, nil
}
