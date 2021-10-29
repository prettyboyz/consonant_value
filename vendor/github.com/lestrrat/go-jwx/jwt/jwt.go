// Package jwt implements JSON Web Tokens as described in https://tools.ietf.org/html/rfc7519
package jwt

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat/go-jwx/internal/emap"
)

// MarshalJSON generates JSON representation of this instant
func (n NumericDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.UTC().Format(numericDateFmt))
}

// UnmarshalJSON parses the JSON representation and initializes this NumericDate
func (n *NumericDate) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t, err := time.Parse(numericDateFmt, s)
	if err != nil {
		return err
	}

	*n = NumericDate{t}
	return nil
}

// NewClaimSet creates a new ClaimSet
func NewClaimSet() *ClaimSet {
	return &ClaimSet{
		EssentialClaims: &EssentialClaims{},
		PrivateClaims:   map[string]interface{}{},
	}
}

// MarshalJSON generates JSON representation of this claim set
func (c *ClaimSet) MarshalJSON() ([]byte, error) {
	// Reverting time back for machines whose time is not perfectly in sync.
	// If client machine's time is in the future according
	// to Google servers, an access token will not be issued.
	now := time.Now().Add(-10 * time.Second)
	if c.IssuedAt == 0 {
		c.IssuedAt = now.Unix()
	}
	if c.Expiration == 0 {
		c.Expiration = now.Add(time.Hour).