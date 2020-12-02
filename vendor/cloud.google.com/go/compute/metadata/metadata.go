// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package metadata provides access to Google Compute Engine (GCE)
// metadata and API service accounts.
//
// This package is a wrapper around the GCE metadata service,
// as documented at https://developers.google.com/compute/docs/metadata.
package metadata // import "cloud.google.com/go/compute/metadata"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const (
	// metadataIP is the documented metadata server IP address.
	metadataIP = "169.254.169.254"

	// metadataHostEnv is the environment variable specifying the
	// GCE metadata hostname.  If empty, the default value of
	// metadataIP ("169.254.169.254") is used instead.
	// This is variable name is not defined by any spec, as far as
	// I know; it was made up for the Go package.
	metadataHostEnv = "GCE_METADATA_HOST"

	userAgent = "gcloud-golang/0.1"
)

type cachedValue struct {
	k    string
	trim bool
	mu   sync.Mutex
	v    string
}

var (
	projID  = &cachedValue{k: "project/project-id", trim: true}
	projNum = &cachedValue{k: "project/numeric-project-id", trim: true}
	instID  = &cachedValue{k: "instance/id", trim: true}
)

var (
	metaClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			ResponseHeaderTimeout: 2 * time.Second,
		},
	}
	subscribeClient = &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
		},
	}
)

// NotDefinedError is returned when requested metadata is not defined.
//
// The underlying string is the suffix after "/computeMetadata/v1/".
//
// This error is not returned if the value is defined to be the empty
// string.
type NotDefinedError string

func (suffix NotDefinedError) Error() string {
	return fmt.Sprintf("metadata: GCE metadata %q not defined", string(suffix))
}

// Get returns a value from the metadata service.
// The suffix is appended to "http://${GCE_METADATA_HOST}/computeMetadata/v1/".
//
// If the GCE_METADATA_HOST environment variable is not defined, a default of
// 169.254.169.254 will be used instead.
//
// If the requested metadata is not defined, the returned error will
// be of type NotDefinedError.
func Get(suffix string) (string, error) {
	val, _, err := getETag(metaClient, suffix)
	return val, err
}

// getETag returns a value from the metadata service as well as the associated
// ETag using the provided client. This func is otherwise equivalent to Get.
func getETag(client *http.Client, suffix string) (value, etag string, err error) {
	// Using a fixed IP makes it very difficult to spoof the metadata service in
	// a container, which is an important use-case for local testing of cloud
	// deployments. To enable spoofing of the metadata service, the environment
	// variable GCE_METADATA_HOST is first inspected to decide where metadata
	// requests shall go.
	host := os.Getenv(metadataHostEnv)
	if host == "" {
		// Using 169.254.169.254 instead of "metadata" here because Go
		// binaries built with the "netgo" tag and without cgo won't
		// know the search suffix for "metadata" is
		// ".google.internal", and this IP address is documented as
		// being stable anyway.
		host = metadataIP
	}
	url := "http://" + host + "/computeMetadata/v1/" + suffix
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Metadata-Flavor", "Google")
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return "", "", NotDefinedError(suffix)
	}
	if res.StatusCode != 200 {
		return "", "", fmt.Errorf("status code %d trying to fetch %s", res.StatusCode, url)
	}
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}
	return string(all), res.Header.Get("Etag"), nil
}

func getTrimmed(suffix string) (s string, err error) {
	s, err = Get(suffix)
	s = strings.TrimSpace(s)
	return
}

func (c *cachedValue) get() (v string, err error) {
	defer c.mu.Unlock()
	c.mu.Lock()
	if c.v != "" {
		return c.v, nil
	}
	if c.trim {
		v, err = getTrimmed(c.k)
	} else {
		v, err = Get(c.k)
	}
	if err == nil {
		c.v = v
	}
	return
}

var (
	onGCEOnce sync.Once
	onGCE     bool
)

// OnGCE reports whether this process is running on Google Compute Engine.
func OnGCE() bool {
	onGCEOnce.Do(initOnGCE)
	return onGCE
}

func initOnGCE() {
	onGCE = testOnGCE()
}

func testOnGCE() bool {
	// The user explicitly said they're on GCE, so trust them.
	if os.Getenv(metadataHostEnv) != "" {
		return true
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resc := make(chan bool, 2)

	// Try two strategies in parallel.
	// See https://github.com/GoogleCloudPlatform/google-cloud-go/issues/194
	go func() {
		req, _ := http.NewRequest("GET", "http://"+metadataIP, nil)
		req.Header.Set("User-Agent", userAgent)
		res, err := ctxht