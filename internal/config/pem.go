package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// PrivateKey represents a parsed RSA private key.
type PrivateKey struct {
	Value *rsa.PrivateKey
}

// UnmarshalText converts an environment variable string to a PrivateKey.
func (k *PrivateKey) UnmarshalText(text []byte) error {
	block, _ := pem.Decode(text)
	if block == nil {
		return fmt.Errorf("no valid PEM data provided")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	k.Value = key

	return nil
}
