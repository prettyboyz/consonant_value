package jwk

import (
	"crypto"
	"crypto/rsa"
	"errors"
	"fmt"
	"math/big"

	"github.com/lestrrat/go-jwx/buffer"
)

// NewRsaPublicKey creates a new JWK using the given key
func NewRsaPublicKey(pk *rsa.PublicKey) (*RsaPublicKey, error) {
	k := &RsaPublicKey{
		EssentialHeader: &EssentialHeader{KeyType: "RSA"},
		N:               buffer.Buffer(pk.N.Bytes()),
		E:               buffer.FromUint(uint64(pk.E)),
	}
	return k, nil
}

// NewRsaPrivateKey creates a new JWK using the given key
func NewRsaPrivateKey(pk *rsa.PrivateKey) (*RsaPrivateKey, error) {
	if len(pk.Primes) < 2 {
		return nil, errors.New("two primes required for RSA private key")
	}

	pub, err := NewRsaPublicKey(&pk.PublicKey)
	if err != nil {
		return nil, err
	}

	k := &RsaPrivateKey{
		RsaPublicKey: pub,
		D:            buffer.Buffer(pk.D.Bytes()),
		P:            buffer.Buffer(pk.Primes[0].Bytes()),
		Q:            buffer.Buffer(pk.Primes[1].Bytes()),
	}

	return k, nil
}

// Materialize returns the RSA public key represented by this JWK
func (k *RsaPublicKey) Materialize() (interface{}, error) {
	return k.PublicKey()
}

// PublicKey creates a new rsa.PublicKey from the data given in the JWK
func (k *RsaPublicKey) PublicKey() (*rsa.PublicKey, error) {
	if k.N.Len() == 0 {
		ret