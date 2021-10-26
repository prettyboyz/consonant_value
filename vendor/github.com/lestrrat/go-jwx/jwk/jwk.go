// Package jwk implements JWK as described in https://tools.ietf.org/html/rfc7517
package jwk

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/lestrrat/go-jwx/internal/emap"
	"github.com/lestrrat/go-jwx/jwa"
)

// FetchFile fetches the local JWK from file, and parses its contents
func FetchFile(jwkpath string) (*Set, error) {
	f, err := os.Open(jwkpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return Parse(buf)
}

// FetchHTTP fetches the remote JWK and parses its contents
func FetchHTTP(jwkurl string) (*Set, error) {
	res, err := http.Get(jwkurl)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch JWK from remote url")
	}

	// XXX Check for maximum length to read?
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return Parse(buf)
}

// Parse parses JWK from the incoming byte buffer.
func Parse(buf []byte) (*Set, error) {
	m := make(map[string]interface{})
	if err := json.Unmarshal(buf, &m); err != nil {
		return nil, err
	}

	// We must change what the underlying structure that gets decoded
	// out of this JSON is based on parameters within the already parsed
	// JSON (m). In order to do this, we have to go through the tedious
	// task of parsing the contents of this map :/
	if _, ok := m["keys"]; ok {
		return constructSet(m)
	}
	k, err := constructKey(m)
	if err != nil {
		return nil, err
	}
	return &Set{Keys: []Key{k}}, nil
}

// ParseString parses JWK from the incoming string.
func ParseString(s string) (*Set, error) {
	return Parse([]byte(s))
}

func constructKey(m map[string]interface{}) (Key, error) {
	kty, ok := m["kty"].(string)
	if !ok {
		return nil, ErrUnsupportedKty
	}

	switch jwa.KeyType(kty) {
	case jwa.RSA:
		if _, ok := m["d"]; ok {
			return constructRsaPrivateKey(m)
		}
		return constructRsaPublicKey(m)
	case jwa.EC:
		if _, ok := m["d"]; ok {
			return constructEcdsaPrivateKey(m)
		}
		return constructEcdsaPublicKey(m)
	case jwa.OctetSeq:
		return constructSymmetricKey(m)
	default:
		return nil, ErrUnsupportedKty
	}
}

func constructEssentialHeader(m map[string]interface{}) (*EssentialHeader, error) {
	r := emap.Hmap(m)
	e := &EssentialHeader{}

	// https://tools.ietf.org/html/rfc7517#section-4.1
	kty, err := r.GetString("kty")
	if err != nil {
		return nil, err
	}
	e.KeyType = jwa.KeyType(kty)

	// https://tools.ietf.org/html/rfc7517#section-4.2
	e.KeyUsage, _ = r.GetString("use")

	// https://tools.ietf.org/html/rfc7517#section-4.3
	if v, err := r.GetStringSlice("key_ops"); err == nil {
		if len(v) > 0 {
			e.KeyOps = make([]KeyOperation, len(v))
			for i, x := range v {
				e.KeyOps[i] = KeyOperation(x)
			}
		}
	}

	// https://tools.ietf.org/html/rfc7517#section-4.4
	e.Algorithm, _ = r.GetString("alg")

	// https://tools.ietf.org/html/rfc7517#section-4.5
	e.KeyID, _ = r.GetString("kid")

	// https://tools.ietf.org/html/rfc7517#section-4.6
	if v, err := r.GetString("x5u"); err == nil {
		u, err := url.Parse(v)
		if err != nil {
			return nil, err
		}
		e.X509Url = u
	}

	// https://tools.ietf.org/html/rfc7517#section-4.7
	if v, err := r.GetStringSlice("x5c"); err == nil {
		e.X509CertChain = v
	}

	return e, nil
}

func constructSymmetricKey(m map[string]interface{}) (*SymmetricKey, error) {
	r := emap.Hmap(m)

	h, err := constructEssentialHeader(m)
	if err != nil {
		return nil, err
	}

	key := &SymmetricKey{EssentialHeader: h}

	k, err := r.GetBuffer("k")
	if err != nil {
		return nil, err
	}
	key.Key = k

	return key, nil
}

func constructEcdsaPublicKey(m map[string]interface{}) (*EcdsaPublicKey, error) {
	e, err := constructEssentialHeader(m)
	if err != nil {
		return nil, err
	}
	r := emap.Hmap(m)

	crvstr, err := r.GetString("crv")
	if err != nil {
		return nil, 