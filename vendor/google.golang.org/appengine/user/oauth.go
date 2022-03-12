// Copyright 2012 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package user

import (
	"golang.org/x/net/context"

	"google.golang.org/appengine/internal"
	pb "google.golang.org/appengine/internal/user"
)

// CurrentOAuth returns the user associated with the OAuth consumer making this
// request. If the OAuth consumer did not make a valid OAuth request, or the
// scopes is non-empty and the current user does not have at least one of t