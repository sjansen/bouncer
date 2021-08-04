package keyring

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"time"
)

type KeyStore interface {
	DescribeJWKs(ctx context.Context) (map[string]time.Time, error)
	PutJWK(ctx context.Context, name, value string) error
	PutSAMLCertificate(ctx context.Context, value string) error
	PutSAMLPrivateKey(ctx context.Context, value string) error
}

type KeyRing struct {
	store KeyStore
}

func New(store KeyStore) *KeyRing {
	return &KeyRing{
		store: store,
	}
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
	keys, err := k.store.DescribeJWKs(ctx)
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
