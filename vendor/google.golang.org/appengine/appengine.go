// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package appengine provides basic functionality for Google App Engine.
//
// For more information on how to write Go apps for Google App Engine, see:
// https://cloud.google.com/appengine/docs/go/
package appengine // import "google.golang.org/appengine"

import (
	"net/http"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	"google.golang.org/appengine/internal"
)

// The gophers party all night; the rabbits provide the beats.

// Main is the principal entry point for an app running in App Engine.
//
// On App Engine Flexible it installs a trivial health checker if one isn't
// already registered, and starts listening on port 8080 (overridden by the
// $PORT environment variable).
//
// See https://cloud.google.com/appengine/docs/flexible/custom-runtimes#hea