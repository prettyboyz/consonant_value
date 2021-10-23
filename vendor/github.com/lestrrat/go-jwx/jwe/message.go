package jwe

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/lestrrat/go-jwx/buffer"
	"github.com/lestrrat/go-jwx/internal/debug"
	"github.com/lestrrat/go-jwx/internal/emap"
	"github.com/lestrrat/go-jwx/jwa"
	"github.com/lestrrat/go-jwx/jwk"
)

// NewRecipient creates a Recipient object
func NewRecipient() *Recipient {
	return &Recipient{
		Header: NewHeader(),
	}
}

// NewHeader creates a new Header object
func NewHeader() *Header {
	return &Header{
		EssentialHeader: &EssentialHeader{},
		PrivateParams:   map[string]interface{}{},
	}
}

// NewEncodedHeader creates a new encoded Header object
func NewEncodedHeader() *EncodedHeader {
	return &EncodedHeader{
		Header: NewHeader(),
	}
}

// Get returns the header key
func (h *Header) Get(key string) (interface{}, error) {
	switch key {
	case "alg":
		return h.Algorithm, nil
	case "apu":
		return h.AgreementPartyUInfo, nil
	case "apv":
		return h.AgreementPartyVInfo, nil
	case "enc":
		return h.ContentEncryption, nil
	case "epk":
		return h.EphemeralPublicKey, nil
	case "cty":
		return h.ContentType, nil
	case "kid":
		return h.KeyID, nil
	case "typ":
		return h.Type, nil
	case "x5t":
		return h.X509CertThumbprint, nil
	case "x5t#256":
		return h.X509CertThumbprintS256, nil
	case "x5c":
		return h.X509CertChain, nil
	case "crit":
		return h.Critical, nil
	case "jku":
		return h.JwkSetURL, nil
	case "x5u":
		return h.X509Url, nil
	default:
		v, ok := h.PrivateParams[key]
		if !ok {
			return nil, errors.New("invalid header name")
		}
		return v, nil
	}
}

// Set sets the value of the given key to the given value. If it's
// one of the known keys, it will be set in EssentialHeader field.
// Otherwise, it is set in PrivateParams field.
func (h *Header) Set(key string, value interface{}) error {
	switch key {
	case "alg":
		var v jwa.KeyEncryptionAlgorithm
		s, ok := value.(string)
		if ok {
			v = jwa.KeyEncryptionAlgorithm(s)
		} else {
			v, ok = value.(jwa.KeyEncryptionAlgorithm)
			if !ok {
				return ErrInvalidHeaderValue
			}
		}
		h.Algorithm = v
	case "apu":
		var v buffer.Buffer
		switch value.(type) {
		case buffer.Buffer:
			v = value.(buffer.Buffer)
		case []byte:
			v = buffer.Buffer(value.([]byte))
		case string:
			v = buffer.Buffer(value.(string))
		default:
			return ErrInvalidHeaderValue
		}
		h.AgreementPartyUInfo = v
	case "apv":
		var v buffer.Buffer
		switch value.(type) {
		case buffer.Buffer:
			v = value.(buffer.Buffer)
		case []byte:
			v = buffer.Buffer(value.([]byte))
		case string:
			v = buffer.Buffer(value.(string))
		default:
			return ErrInvalidHeaderValue
		}
		h.AgreementPartyVInfo = v
	case "enc":
		var v jwa.ContentEncryptionAlgorithm
		s, ok := value.(string)
		if ok {
			v = jwa.ContentEncryptionAlgorithm(s)
		} else {
			v, ok = value.(jwa.ContentEncryptionAlgorithm)
			if !ok {
				return ErrInvalidHeaderValue
			}
		}
		h.ContentEncryption = v
	case "cty":
		v, ok := value.(string)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.ContentType = v
	case "epk":
		v, ok := value.(*jwk.EcdsaPublicKey)
		if !ok {
			return ErrInvalidHeaderValue
		}
		h.EphemeralPublicKey = v
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

	for k, v := range h2.PrivateParams