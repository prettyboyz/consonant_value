// Package jws implements the digital signature on JSON based data
// structures as described in https://tools.ietf.org/html/rfc7515
//
// If you do not care about the details, the only things that you
// would need to use are the following functions:
//
//     jws.Sign(payload, algorithm, key)
//     jws.Verify(encodedjws, algorithm, key)
//
// To sign, simply use `jws.Sign`. `payload` is a []byte buffer that
// contains whatever data you want to sign. `alg` is one of the
// jwa.SignatureAlgorithm constants from package jwa. For RSA and
// ECDSA family of algorithms, you will need to prepare a private key.
// For HMAC family, you just need a []byte value. The `jws.Sign`
// function will return the encoded JWS message on success.
//
// To verify, use `jws.Verify`. It will parse the `encodedjws` buffer
// and verify the result using `algorithm` and `key`. Upon successful
// verification, the original payload is returned, so you can work on it.
package jws

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/lestrrat/go-jwx/buffer"
	"github.com/lestrrat/go-jwx/internal/debug"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwk"
)

// Sign is a short way to generate a JWS in compact serialization
// for a given payload. If you need more control over the signature
// generation process, you should manually create signers and tweak
// the message.
//
//
func Sign(payload []byte, alg jwa.SignatureAlgorithm, key interface{}, hdrs ...*Header) ([]byte, error) {
	var err error
	var signer PayloadSigner
	switch alg {
	case jwa.RS256, jwa.RS384, jwa.RS512, jwa.PS256, jwa.PS384, jwa.PS512:
		privkey, ok := key.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("invalid private key: *rsa.PrivateKey required")
		}

		signer, err = NewRsaSign(alg, privkey)
		if err != nil {
			return nil, err
		}
	case jwa.HS256, jwa.HS384, jwa.HS512:
		sharedkey, ok := key.([]byte)
		if !ok {
			return nil, errors.New("invalid private key: []byte required")
		}

		signer, err = NewHmacSign(alg, sharedkey)
		if err != nil {
			return nil, err
		}
	case jwa.ES256, jwa.ES384, jwa.ES512:
		privkey, ok := key.(*ecdsa.PrivateKey)
		if !ok {
			return nil, errors.New("invalid private key: *ecdsa.PrivateKey required")
		}

		signer, err = NewEcdsaSign(alg, privkey)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportedAlgorithm
	}

	if len(hdrs) > 0 {
		if pubhdr := hdrs[0]; pubhdr != nil {
			h, err := signer.PublicHeaders().Merge(pubhdr)
			if err != nil {
				return nil, err
			}
			signer.SetPublicHeaders(h)
		}
	}

	if len(hdrs) > 1 {
		if protectedhdr := hdrs[0]; protectedhdr != nil {
			h, err := signer.ProtectedHeaders().Merge(protectedhdr)
			if err != nil {
				return nil, err
			}
			signer.SetProtectedHeaders(h)
		}
	}

	multisigner := NewMultiSign()
	multisigner.AddSigner(signer)
	msg, err := multisigner.Sign(payload)
	if err != nil {
		return nil, err
	}

	return CompactSerialize{}.Serialize(msg)
}

// Verify checks if the given JWS message is verifiable using `alg` and `key`.
// If the verification is successful, `err` is nil, and the content of the
// payload that was signed is returned. If you need more fine-grained
// control of the verification process, manually call `Parse`, generate a
// verifier, and call `Verify` on the parsed JWS message object.
func Verify(buf []byte, alg jwa.SignatureAlgorithm, key interface{}) ([]byte, error) {
	if debug.Enabled {
		debug.Printf("jws.Verify\n")
	}
	msg, err := Parse(buf)
	if err != nil {
		return nil, err
	}

	var verifier Verifier
	switch alg {
	case jwa.RS256, jwa.RS384, jwa.RS512, jwa.PS256, jwa.PS384, jwa.PS512:
		pubkey, ok := key.(*rsa.PublicKey)
		if !ok {
			return nil, errors.New("invalid key: *rsa.PublicKey required")
		}

		rsaverify, err := NewRsaVerify(alg, pubkey)
		if err != nil {
			return nil, err
		}
		verifier = rsaverify
	case jwa.HS256, jwa.HS384, jwa.HS512:
		sharedkey, ok := key.([]byte)
		if !ok {
			return nil, errors.New("invalid key: []byte required")
		}

		hmacverify, err := NewHmacVerify(alg, sharedkey)
		if err != nil {
			return nil, err
		}
		verifier = hmacverify
	case jwa.ES256, jwa.ES384, jwa.ES512:
		pubkey, ok := key.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.New("invalid key: *ecdsa.PublicKey required")
		}

		ecdsaverify, err := NewEcdsaVerify(alg, pubkey)
		