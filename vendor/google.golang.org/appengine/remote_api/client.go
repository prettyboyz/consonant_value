// Copyright 2013 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package remote_api

// This file provides the client for connecting remotely to a user's production
// application.

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	"google.golang.org/appengine/internal"
	pb "google.golang.org/appengine/internal/remote_api"
)

// NewRemoteContext returns a context that gives access to the production
// APIs for the application at the given host. All communication will be
// performed over SSL unless the host is localhost.
func NewRemoteContext(host string, client *http.Client) (context.Context, error) {
	// Add an appcfg header to outgoing requests.
	t := client.Transport
	if t == nil {
		t = http.DefaultTransport
	}
	client.Transport = &headerAddingRoundTripper{t}

	url := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/_ah/remote_api",
	}
	if host == "localhost" || strings.HasPrefix(host, "localhost:") {
		url.Scheme = "http"
	}
	u := url.String()
	appID, err := getAppID(client, u)
	if err != nil {
		return nil, fmt.Errorf("unable to contact server: %v", err)
	}
	rc := &remoteContext{
		client: client,
		url:    u,
	}
	ctx := internal.WithCallOverride(context.Background(), rc.call)
	ctx = internal.WithLogOverride(ctx, rc.logf)
	ctx = internal.WithAppIDOverride(ctx, appID)
	return ctx, nil
}

type remoteContext struct {
	client *http.Client
	url    string
}

var logLevels = map[int64]string{
	0: "DEBUG",
	1: "INFO",
	2: "WARNING",
	3: "ERROR",
	4: "CRITICAL",
}

func (c *remoteContext) logf(level int64, format string, args ...interface{}) {
	log.Printf(logLevels[level]+": "+format, args...)
}

func (c *remoteContext) call(ctx context.Context, service, method string, in, out proto.Message) error {
	req, err := proto.Marshal(in)
	if err != nil {
		return fmt.Errorf("error marshalling request: %v", err)
	}

	remReq := &pb.Request{
		ServiceName: proto.String(service),
		Method:      proto.String(method),
		Request:     req,
		// NOTE(djd): RequestId is unused in the server.
	}

	req, err = proto.Marshal(remReq)
	if err != nil {
		return fmt.Errorf("proto.Marshal: %v", err)
	}

	// TODO(djd): Respect ctx.Deadline()?
	resp, err := c.client.Post(c.url, "application/octet-stream", bytes.NewReader(req))
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response %d; body: %q", resp.StatusCode, body)
	}
	if err != nil {
		return fmt.Errorf("failed reading response: %v", err)
	}
	remResp := &pb.Response{}
	if err := proto.Unmarshal(body, remResp); err != nil {
		return fmt.Errorf("error unmarshalling response: %v", err)
	}

	if ae := remResp.GetApplicationError(); ae != nil {
		return &internal.APIError{
			Code:    ae.GetCode(),
			Detail:  ae.GetDetail(),
			Service: service,
		}
	}

	if remResp.Response == nil {
		return fmt