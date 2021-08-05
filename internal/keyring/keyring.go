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

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
)

type KeyStore interface {
	DescribeJWKSet(ctx context.Context) (map[string]time.Time, error)
	GetJWKSet(ctx context.Context) (map[string]string, map[string]time.Time, error)
	PutJWK(ctx context.Context, name, value string) error
	PutSAMLCertificate(ctx context.Context, value string) error
	PutSAMLPrivateKey(ctx context.Context, value string) error
}

type KeyRing struct {
	sync.RWMutex
	json  []byte
	store KeyStore
	jwt   struct {
		key   *rsa.PrivateKey
		kid   string
		mtime time.Time
	}
}

func New(store KeyStore) *KeyRing {
	return &KeyRing{
		json:  []byte("{}"),
		store: store,
	}
}

// JWKSetAsJSON returns the current JSON Web Key Set.
func (k *KeyRing) JWKSetAsJSON() []byte {
	k.RLock()
	defer k.RUnlock()
	return k.json
}

// NewJWT returns a few JWT for subject.
func (k *KeyRing) NewJWT(subj string) ([]byte, error) {
	now := time.Now().UTC()
	h := jws.NewHeaders()
	_ = h.Set(jws.KeyIDKey, k.jwt.kid)
	t := jwt.New()
	_ = t.Set(jwt.SubjectKey, subj)
	_ = t.Set(jwt.AudienceKey, "bouncer-authz")
	_ = t.Set(jwt.IssuerKey, "bouncer")
	_ = t.Set(jwt.IssuedAtKey, now)
	_ = t.Set(jwt.ExpirationKey, now.Add(1*time.Hour))
	return jwt.Sign(
		t, jwa.RS256, k.jwt.key,
		jwt.WithJwsHeaders(h),
	)
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

// WatchJWKSet polls for updated JSON Web Keys.
func (k *KeyRing) WatchJWKSet(ctx context.Context) error {
	set, err := k.getjwkset(ctx)
	if err != nil {
		return err
	}

	jwks, err := json.Marshal(set)
	if err != nil {
		return err
	}

	k.Lock()
	defer k.Unlock()
	k.json = jwks

	return nil
}

func (k *KeyRing) getjwkset(ctx context.Context) (jwk.Set, error) {
	keys, mtimes, err := k.store.GetJWKSet(ctx)
	if err != nil {
		return nil, err
	}

	set := jwk.NewSet()
	for name, value := range keys {
		block, _ := pem.Decode([]byte(value))
		if block == nil {
			return nil, fmt.Errorf("invalid PEM data: %s", name)
		}

		parsed, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		kid, key, err := genjwk(parsed)
		if err != nil {
			return nil, err
		}
		set.Add(key)

		if mtime, ok := mtimes[name]; ok && k.jwt.mtime.Before(mtime) {
			k.jwt.kid = kid
			k.jwt.key = parsed
			k.jwt.mtime = mtime
		}
	}

	return set, nil
}

func genjwk(key *rsa.PrivateKey) (string, jwk.Key, error) {
	k, err := jwk.New(&key.PublicKey)
	if err != nil {
		return "", nil, err
	}

	err = jwk.AssignKeyID(k)
	if err != nil {
		return "", nil, err
	}

	err = k.Set(jwk.KeyUsageKey, "sig")
	if err != nil {
		return "", nil, err
	}

	kid, _ := k.Get(jwk.KeyIDKey)
	return kid.(string), k, nil
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
