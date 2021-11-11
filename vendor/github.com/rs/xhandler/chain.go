
package xhandler

import (
	"net/http"

	"golang.org/x/net/context"
)

// Chain is a helper for chaining middleware handlers together for easier
// management.
type Chain []func(next HandlerC) HandlerC

// Add appends a variable number of additional middleware handlers
// to the middleware chain. Middleware handlers can either be
// context-aware or non-context aware handlers with the appropriate
// function signatures.
func (c *Chain) Add(f ...interface{}) {
	for _, h := range f {
		switch v := h.(type) {
		case func(http.Handler) http.Handler:
			c.Use(v)
		case func(HandlerC) HandlerC:
			c.UseC(v)
		default:
			panic("Adding invalid handler to the middleware chain")
		}
	}
}

// With creates a new middleware chain from an existing chain,
// extending it with additional middleware. Middleware handlers
// can either be context-aware or non-context aware handlers
// with the appropriate function signatures.
func (c *Chain) With(f ...interface{}) *Chain {
	n := make(Chain, len(*c))
	copy(n, *c)
	n.Add(f...)
	return &n
}

// UseC appends a context-aware handler to the middleware chain.
func (c *Chain) UseC(f func(next HandlerC) HandlerC) {
	*c = append(*c, f)
}

// Use appends a standard http.Handler to the middleware chain without
// losing track of the context when inserted between two context aware handlers.
//