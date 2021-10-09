// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mux

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Route stores information to match a request and build URLs.
type Route struct {
	// Parent where the route was registered (a Router).
	parent parentRoute
	// Request handler for the route.
	handler http.Handler
	// List of matchers.
	matchers []matcher
	// Manager for the variables from host and path.
	regexp *routeRegexpGroup
	// If true, when the path pattern is "/path/", accessing "/path" will
	// redirect to the former and vice versa.
	strictSlash bool
	// If true, when the path pattern is "/path//to", accessing "/path//to"
	// will not redirect
	skipClean bool
	// If true, "/path/foo%2Fbar/to" will match the path "/path/{var}/to"
	useEncodedPath bool
	// If true, this route never matches: it is only used to build URLs.
	buildOnly bool
	// The name used to build URLs.
	name string
	// Error resulted from building a route.
	err error

	buildVarsFunc BuildVarsFunc
}

func (r *Route) SkipClean() bool {
	return r.skipClean
}

// Match matches the route against the request.
func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
	if r.buildOnly || r.err != nil {
		return false
	}
	// Match everything.
	for _, m := range r.matchers {
		if matched := m.Match(req, match); !matched {
			return false
		}
	}
	// Yay, we have a match. Let's collect some inf