package jws

import (
	"encoding/json"
	"errors"
	"net/url"

	"github.com/lestrrat/go-jwx/buffer"
	"github.com/lestrrat/go-jwx/internal/emap"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwk"
)

// NewHeader creates a new Header
func NewHeader() *Header {
	return &Header{
		EssentialHeader: &EssentialHeader{},
		PrivateParams:   map[string]interface{}{},
	}
}

// Set sets the value of the given key to the given value. If it's
// one of the known keys, it will be set in EssentialHeader field.
// Otherwise, it is set in PrivateParams field.
func (h *Header) Set(key string, value interface{}) error {
	switch key {
	case "alg":
		var v jwa.SignatureAlgorithm
		s, ok := value.(string)
		if ok {
			v = jwa.SignatureAlgorithm(s)
		} else {
			v, ok = value.(jwa.SignatureAlgorithm)
			if !ok {
				return ErrInvalidHeaderValue
			}
		}
		h.Algorithm = v
	case "cty":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.ContentType = v
	case "kid":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.KeyID = v
	case "typ":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.Type = v
	case "x5t":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.X509CertThumbprint = v
	case "x5t#256":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.X509CertThumbprintS256 = v
	case "x5c":
		v, ok := value.([]string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.X509CertChain = v
	case "crit":
		v, ok := value.([]string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.Critical = v
	case "jwk":
		v, ok := value.(jwk.Key)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.Jwk = v
	case "jku":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		u, err := url.Parse(v)
		if err != nil {
			return ErrInvalidHeaderValue
		}
		h.JwkSetURL = u
	case "x5u":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		u, err := url.Parse(v)
		if err != nil {
			return ErrInvalidHeaderValue
		}
		h.X509Url = u
	default:
		h.PrivateParams[key] = value
	}
	return nil
}

// Merge merges the current header with another.
func (h *Header) Merge(h2 *Header) (*Header, error) {
	if h2 == nil {
		return nil, errors.New("merge target is nil")
	}

	h3 := NewHeader()
	if err := h3.Copy(h); err != nil {
		return nil, err
	}

	h3.EssentialHeader.Merge(h2.EssentialHeader)

	for k, v := range h2.PrivateParams {
		h3.PrivateParams[k] = v
	}

	return h3, nil
}

// Merge merges the current header with another.
func (h *EssentialHeader) Merge(h2 *EssentialHeader) {
	if h2.Algorithm != "" {
		h.Algorithm = h2.Algorithm
	}

	if h2.ContentType != "" {
		h.ContentType = h2.ContentType
	}

	if h2.Jwk != nil {
		h.Jwk = h2.Jwk
	}

	if h2.JwkSetURL != nil {
		h.JwkSetURL = h2.JwkSetURL
	}

	if h2.KeyID != "" {
		h.KeyID = h2.KeyID
	}

	if h2.Type != "" {
		h.Type = h2.Type
	}

	if h2.X509Url != nil {
		h.X509Url = h2.X509Url
	}

	if h2.X509CertChain != nil {
		h.X509CertChain = h2.X509CertChain
	}

	if h2.X509CertThumbprint != "" {
		h.X509CertThumbprint = h2.X509CertThumbprint
	}

	if h2.X509CertThumbprintS256 != "" {
		h.X509CertThumbprintS256 = h2.X509CertThumbprintS256
	}
}

// Copy copies the other he