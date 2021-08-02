package config

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRekey(t *testing.T) {
	before := time.Now()
	after := before.Add(1 * time.Minute)
	described := map[string]struct{}{
		"JWK1": struct{}{},
		"JWK2": struct{}{},
	}

	for _, tc := range []struct {
		label   string
		mtimes  map[string]time.Time
		updated map[string]struct{}
	}{{
		label: "both uninitialized",
		updated: map[string]struct{}{
			"JWK1": struct{}{},
			"JWK2": struct{}{},
		},
	}, {
		label: "key1 uninitialized",
		mtimes: map[string]time.Time{
			"JWK2": after,
		},
		updated: map[string]struct{}{
			"JWK1": struct{}{},
		},
	}, {
		label: "key2 uninitialized",
		mtimes: map[string]time.Time{
			"JWK1": before,
		},
		updated: map[string]struct{}{
			"JWK2": struct{}{},
		},
	}, {
		label: "jwk1 older",
		mtimes: map[string]time.Time{
			"JWK1": before,
			"JWK2": after,
		},
		updated: map[string]struct{}{
			"JWK1": struct{}{},
		},
	}, {
		label: "jwk2 older",
		mtimes: map[string]time.Time{
			"JWK1": after,
			"JWK2": before,
		},
		updated: map[string]struct{}{
			"JWK2": struct{}{},
		},
	}} {
		t.Run(tc.label, func(t *testing.T) {
			require := require.New(t)
			tc := tc

			c := &mockSSMClient{mtimes: tc.mtimes}
			err := rekey(context.TODO(), c)
			require.NoError(err)
			require.Equal(described, c.described)
			require.Equal(tc.updated, c.updated)
		})
	}
}

type mockSSMClient struct {
	ssmClient
	described map[string]struct{}
	updated   map[string]struct{}
	mtimes    map[string]time.Time
}

func (c *mockSSMClient) DescribeParameters(ctx context.Context, names ...string) (map[string]time.Time, error) {
	if c.described == nil {
		c.described = map[string]struct{}{}
	}
	for _, name := range names {
		c.described[name] = struct{}{}
	}
	return c.mtimes, nil
}

func (c *mockSSMClient) PutParameter(ctx context.Context, key, value string, encrypt bool) error {
	if c.updated == nil {
		c.updated = map[string]struct{}{
			key: struct{}{},
		}
	} else {
		c.updated[key] = struct{}{}
	}
	return nil
}
