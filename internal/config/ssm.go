package config

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/vrischmann/envconfig"

	"github.com/sjansen/bouncer/internal/aws"
)

type ssmConfig struct {
	Prefix string `envconfig:"BOUNCER_SSM_PREFIX"`
}

type ssmClient interface {
	DescribeParameters(ctx context.Context, names ...string) (map[string]time.Time, error)
	PutParameter(ctx context.Context, key, value string, encrypt bool) error
}

func Rekey(ctx context.Context) error {
	cfg := &ssmConfig{}
	if err := envconfig.Init(&cfg); err != nil {
		return err
	}

	aws, err := aws.New(ctx)
	if err != nil {
		return err
	}

	ssm, err := aws.NewSSMClient(cfg.Prefix)
	if err != nil {
		return err
	}

	return rekey(ctx, ssm)
}

func genkey() (jwk.Key, error) {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("failed to generate new RSA privatre key: %s\n", err)
		return nil, err
	}

	key, err := jwk.New(raw)
	if err != nil {
		fmt.Printf("failed to create symmetric key: %s\n", err)
		return nil, err
	}

	return key, nil
}

func rekey(ctx context.Context, ssm ssmClient) error {
	keys, err := ssm.DescribeParameters(ctx, "JWK1", "JWK2")
	if err != nil {
		return err
	}

	mtime1, ok1 := keys["JWK1"]
	mtime2, ok2 := keys["JWK2"]
	if !ok1 || ok2 && mtime1.Before(mtime2) {
		key, err := genkey()
		if err != nil {
			return err
		}
		buf, err := json.Marshal(key)
		if err != nil {
			return err
		}
		err = ssm.PutParameter(ctx, "JWK1", string(buf), true)
		if err != nil {
			return err
		}
	}
	if !ok2 || ok1 && !mtime1.Before(mtime2) {
		key, err := genkey()
		if err != nil {
			return err
		}
		buf, err := json.Marshal(key)
		if err != nil {
			return err
		}
		err = ssm.PutParameter(ctx, "JWK2", string(buf), true)
		if err != nil {
			return err
		}
	}

	return nil
}
