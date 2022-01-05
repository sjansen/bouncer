package keyring

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRekeyJWKs(t *testing.T) {
	before := time.Now()
	after := before.Add(1 * time.Minute)

	for _, tc := range []struct {
		label     string
		described bool
		mtimes    map[string]time.Time
		updated   map[string]struct{}
	}{{
		label: "both uninitialized",
		updated: map[string]struct{}{
			"JWK1": {},
			"JWK2": {},
		},
	}, {
		label: "key1 uninitialized",
		mtimes: map[string]time.Time{
			"JWK2": after,
		},
		updated: map[string]struct{}{
			"JWK1": {},
		},
	}, {
		label: "key2 uninitialized",
		mtimes: map[string]time.Time{
			"JWK1": before,
		},
		updated: map[string]struct{}{
			"JWK2": {},
		},
	}, {
		label: "jwk1 older",
		mtimes: map[string]time.Time{
			"JWK1": before,
			"JWK2": after,
		},
		updated: map[string]struct{}{
			"JWK1": {},
		},
	}, {
		label: "jwk2 older",
		mtimes: map[string]time.Time{
			"JWK1": after,
			"JWK2": before,
		},
		updated: map[string]struct{}{
			"JWK2": {},
		},
	}} {
		t.Run(tc.label, func(t *testing.T) {
			require := require.New(t)
			tc := tc

			store := &mockKeyStore{mtimes: tc.mtimes}
			keyring := New(store)
			err := keyring.RotateJWKs(context.TODO())
			require.NoError(err)
			require.Equal(true, store.described)
			require.Equal(tc.updated, store.updated)
		})
	}
}

type mockKeyStore struct {
	KeyStore
	described bool
	mtimes    map[string]time.Time
	updated   map[string]struct{}
}

func (c *mockKeyStore) DescribeJWKSet(ctx context.Context) (map[string]time.Time, error) {
	c.described = true
	return c.mtimes, nil
}

func (c *mockKeyStore) putParameter(key string) error {
	if c.updated == nil {
		c.updated = map[string]struct{}{
			key: {},
		}
	} else {
		c.updated[key] = struct{}{}
	}
	return nil
}

func (c *mockKeyStore) PutJWK(ctx context.Context, name, value string) error {
	return c.putParameter(name)
}
