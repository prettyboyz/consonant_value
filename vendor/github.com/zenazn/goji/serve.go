// +build !appengine

package goji

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/zenazn/goji/bind"
	"github.com/zenazn/goji/graceful"
)

func init() {
	bind.WithFlag()
	if fl := log.Flags(); fl&log.Ltime != 0 {
		log.SetFlags(fl | log.Lmicroseconds)
	}
	graceful.DoubleKickWindow(2 * time.Second)
}

// Serve starts Goji using reasonable defaults.
func Serve() {
	if !flag.Parsed() {
		flag.Parse()
	}

	ServeListener(bind.Default())
}

// Like Serve, but enables TLS using the given config.
func ServeTLS(config *tls.Config) {
	if !flag.Parsed() {
		flag.Parse()
	}

	ServeListener(tls.NewListener(bind.Default(), config))
}

// Like Serve, but runs Goji on top of an arbitra