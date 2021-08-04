package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

// Load reads settings from the environment and AWS SSM.
func Load(ctx context.Context, config interface{}) error {
	ssm, err := newSSMClient(ctx)
	if err != nil {
		return err
	}

	m := &ssmMutator{ssm}
	return load(ctx, config, envconfig.OsLookuper(), m.resolve)
}

func load(ctx context.Context, config interface{}, l envconfig.Lookuper, fns ...envconfig.MutatorFunc) error {
	err := envconfig.ProcessWith(ctx, config,
		envconfig.PrefixLookuper("BOUNCER_", l),
		fns...,
	)
	return err
}
