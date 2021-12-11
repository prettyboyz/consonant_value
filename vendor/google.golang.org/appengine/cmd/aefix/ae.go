// Copyright 2016 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package main

import (
	"go/ast"
	"path"
	"strconv"
	"strings"
)

const (
	ctxPackage = "golang.org/x/net/context"

	newPackageBase = "google.golang.org/"
	stutterPackage = false
)

func init() {
	register(fix{
		"ae",
		"2016-04-15",
		aeFn,
		`Update old App Engine APIs to new App Engine APIs`,
	})
}

// logMethod is the set of methods on appengine.Context used for logging.
var logMethod = map[string]bool{
	"Debugf":    true,
	"Infof":     true,
	"Warningf":  true,
	"Errorf":    true,
	"Criticalf": true,
}

// mapPackage turns "appengine" into "google.golang.org/appengine", etc.
func mapPackage(s string) string {
	if stutterPackage {
		s += "/" + path.Base(s)
	}
	return newPackageBase + s
}

func aeFn(f *ast.File) bool {
	// During the walk, we track the last thing seen that looks like
	// an appeng