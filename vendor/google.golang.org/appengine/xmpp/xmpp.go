// Copyright 2011 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

/*
Package xmpp provides the means to send and receive instant messages
to and from users of XMPP-compatible services.

To send a message,
	m := &xmpp.Message{
		To:   []string{"kaylee@example.com"},
		Body: `Hi! How's the carrot?`,
	}
	err := m.Send(c)

To receive messages,
	func init() {
		xmpp.Handle(handleChat)
	}

	func handleChat(c context.Context, m *xmpp.Message) {
		// ...
	}
*/
package xmpp // import "google.golang.org/appengine/xmpp"

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/internal"
	pb "google.golang.org/appengine/internal/xmpp"
)

// Message represents an incoming chat message.
type Message struct {
	// Sender is the JID of the sender.
	// Optional for outgoing messages.
	Sender string

	// To is the intended recipients of the message.
	// Incoming messages will have exactly one element.
	To []string

	// Body is the body of the message.
	Body string

	// Type is the message type, per RFC 3921.
	// It defaults to "chat".
	Type string

	// RawXML is whether the body contains raw XML.
	RawXML bool
}

// Presence represents an outgoing presence update.
type Presence struct {
	// Sender is the JID (optional).
	Sender string

	// The intended recipient of the presence update.
	To string

	// Type, per RFC 3921 (optional). Defaults to "available".
	Type string

	// State of presence (optional).
	// Valid values: "away", "chat", "xa", "dnd" (RFC 3921).
	State string

	// Free text status message (optional).
	Status string
}

var (
	ErrPresenceUnavailable = errors.New("xmpp: presence unavailable")
	ErrInvalidJID          = errors.New("xmpp: invalid JID")
)

// Handle arranges for f to be called for incoming XMPP messages.
// Only messages of type "chat" or "normal" will be handled.
func Handle(f func(c context.Context, m *Message)) {
	http.HandleFunc("/_ah/xmpp/message/chat/", func(_ http.ResponseWriter, r *http.Request) {
		f(appengine.NewContext(r), &Message{
			Sender: r.FormValue("from"),
			To:     []string{r.FormValue("to")},
			Body:   r.FormValue("body"),
		})
	