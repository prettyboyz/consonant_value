// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

// Package internal provides support for package appengine.
//
// Programs should not use this package directly. Its API is not stable.
// Use packages appengine and appengine/* instead.
package internal

import (
	"fmt"

	"github.com/golang/protobuf/proto"

	remotepb "google.golang.org/appengine/internal/remote_api"
)

// errorCodeMaps is a map of service name to the error code map for the service.
var errorCodeMaps = make(map[string]map[int32]string)

// RegisterErrorCodeMap is called from API implementations to register their
// error code map. This should only be called from init functions.
func RegisterErrorCodeMap(service string, m map[int32]string) {
	errorCodeMaps[service] = m
}

type timeoutCodeKey struct {
	service string
	code    int32
}

// timeoutCodes is the set of service+code pairs that represent timeouts.
var timeoutCodes = make(map[timeoutCodeKey]bool)

func RegisterTimeoutErrorCode(service string, code int32) {
	timeoutCodes[timeoutCodeKey{service, code}] = true
}

// APIError is the typ