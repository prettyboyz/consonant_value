// Copyright 2015 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// +build appengine

package user

import (
	"appengine/user"

	"golang.org/x/net/context"

	"google.golang.org/appengine/internal"
)

func Current(ctx context.Context) *User {
	c, err := in