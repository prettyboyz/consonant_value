
/*
Package web provides a fast and flexible middleware stack and mux.

This package attempts to solve three problems that net/http does not. First, it
allows you to specify flexible patterns, including routes with named parameters
and regular expressions. Second, it allows you to write reconfigurable
middleware stacks. And finally, it allows you to attach additional context to
requests, in a manner that can be manipulated by both compliant middleware and
handlers.
*/
package web

import (
	"net/http"
)

/*
C is a request-local context object which is threaded through all compliant
middleware layers and given to the final request handler.
*/
type C struct {
	// URLParams is a map of variables extracted from the URL (typically
	// from the path portion) during routing. See the documentation for the
	// URL Pattern you are using (or the documentation for PatternType for
	// the case of standard pattern types) for more information about how
	// variables are extracted and named.
	URLParams map[string]string
	// Env is a free-form environment for storing request-local data. Keys
	// may be arbitrary types that support equality, however package-private
	// types with type-safe accessors provide a convenient way for packages
	// to mediate access to their request-local data.
	Env map[interface{}]interface{}
}

// Handler is similar to net/http's http.Handler, but also accepts a Goji
// context object.
type Handler interface {
	ServeHTTPC(C, http.ResponseWriter, *http.Request)
}

// HandlerFunc is similar to net/http's http.HandlerFunc, but supports a context
// object. Implements both http.Handler and Handler.
type HandlerFunc func(C, http.ResponseWriter, *http.Request)

// ServeHTTP implements http.Handler, allowing HandlerFunc's to be used with
// net/http and other compliant routers. When used in this way, the underlying