package web

import (
	"log"
	"net/http"
	"regexp"
)

// A Pattern determines whether or not a given request matches some criteria.
// They are often used in routes, which are essentially (pattern, methodSet,
// handler) tuples. If the method and pattern match, the given handler is used.
//
// Built-in implementations of this interface are used to implement reg