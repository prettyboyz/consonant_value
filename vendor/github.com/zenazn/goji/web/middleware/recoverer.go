package middleware

import (
	"bytes"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/zenazn/goji/web"
)

// Recoverer is a middleware that recovers from panics, logs the pan