package config

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/lestrrat-go/jwx/jwk"
)

// RekeySAML creates or updates the SAML key pair.
func RekeySAML(ctx context.Context) error {
	ssm, err := newSSMClient(ctx)
	if err != nil {
		return err
	}

	return rekeySAML(ctx, ssm)
}

// RotateJWKs creates any missing JSON Web Keys, or replaces the oldest if none are missing.
func RotateJWKs(ctx context.Context) error {
	ssm, err := newSSMClient(ctx)
	if err != nil {
		return err
	}

	return rotateJWKs(ctx, ssm)
}

func genjwk() (jwk.Key, error) {
	raw, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	key, err := jwk.New(raw)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func rekeySAML(ctx context.Context, ssm ssmClient) error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	buf := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})

	err = ssm.PutParameter(ctx, "SAML_PRIVATE_KEY", string(buf), true)
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

	buf = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	})

	err = ssm.PutParameter(ctx, "SAML_CERTIFICATE", string(buf), false)
	if err != nil {
		return err
	}
	return nil
}

func rotateJWKs(ctx context.Context, ssm ssmClient) error {
	keys, err := ssm.DescribeParameters(ctx, "JWK1", "JWK2")
	if err != nil {
		return err
	}

	mtime1, ok1 := keys["JWK1"]
	mtime2, ok2 := keys["JWK2"]
	if !ok1 || ok2 && mtime1.Before(mtime2) {
		key, err := genjwk()
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
		key, err := genjwk()
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
