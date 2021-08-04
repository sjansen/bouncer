package keyring

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
)

type KeyStore interface {
	DescribeJWKSet(ctx context.Context) (map[string]time.Time, error)
	GetJWKSet(ctx context.Context) (map[string]string, error)
	PutJWK(ctx context.Context, name, value string) error
	PutSAMLCertificate(ctx context.Context, value string) error
	PutSAMLPrivateKey(ctx context.Context, value string) error
}

type KeyRing struct {
	sync.RWMutex
	jwks  []byte
	store KeyStore
}

func New(store KeyStore) *KeyRing {
	return &KeyRing{
		jwks:  []byte("{}"),
		store: store,
	}
}

// JWKSetAsJSON returns the current JSON Web Key Set.
func (k *KeyRing) JWKSetAsJSON() []byte {
	k.RLock()
	defer k.RUnlock()
	return k.jwks
}

// RekeySAML creates or updates the SAML key pair.
func (k *KeyRing) RekeySAML(ctx context.Context) error {
	key, encoded, err := genkey()
	if err != nil {
		return err
	}

	err = k.store.PutSAMLPrivateKey(ctx, encoded)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(now.Unix()),
		NotBefore:    now,
		NotAfter:     now.AddDate(10, 0, 0),
	}
	cert, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return err
	}

	buf := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	return k.store.PutSAMLCertificate(ctx, string(buf))
}

// RotateJWKs creates any missing JSON Web Keys, or replaces the oldest if none are missing.
func (k *KeyRing) RotateJWKs(ctx context.Context) error {
	keys, err := k.store.DescribeJWKSet(ctx)
	if err != nil {
		return err
	}

	mtime1, ok1 := keys["JWK1"]
	mtime2, ok2 := keys["JWK2"]
	if !ok1 || ok2 && mtime1.Before(mtime2) {
		_, encoded, err := genkey()
		if err != nil {
			return err
		}
		err = k.store.PutJWK(ctx, "JWK1", encoded)
		if err != nil {
			return err
		}
	}
	if !ok2 || ok1 && !mtime1.Before(mtime2) {
		_, encoded, err := genkey()
		if err != nil {
			return err
		}
		err = k.store.PutJWK(ctx, "JWK2", encoded)
		if err != nil {
			return err
		}
	}

	return nil
}

func genkey() (*rsa.PrivateKey, string, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, "", err
	}

	buf := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	return key, string(buf), nil
}

// WatchJWKSet polls for updates to the JSON Web Key Set.
func (k *KeyRing) WatchJWKSet(ctx context.Context) error {
	set, err := getkeys(ctx, k.store)
	if err != nil {
		return err
	}

	jwks, err := json.Marshal(set)
	if err != nil {
		return err
	}

	k.Lock()
	defer k.Unlock()
	k.jwks = jwks

	return nil
}

func getkeys(ctx context.Context, store KeyStore) (jwk.Set, error) {
	keys, err := store.GetJWKSet(ctx)
	if err != nil {
		return nil, err
	}

	set := jwk.NewSet()
	for k, v := range keys {
		block, _ := pem.Decode([]byte(v))
		if block == nil {
			return nil, fmt.Errorf("invalid PEM data: %s", k)
		}

		parsed, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		key, err := jwk.New(&parsed.PublicKey)
		if err != nil {
			return nil, err
		}

		_ = key.Set(jwk.KeyUsageKey, "sig")
		set.Add(key)
	}

	return set, nil
}
