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
		Priv