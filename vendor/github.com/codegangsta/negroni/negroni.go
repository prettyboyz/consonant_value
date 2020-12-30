package negroni

import (
	"log"
	"net/http"
	"os"
)

const (
	// DefaultAddress is used if no other is specified.
	DefaultAddress = ":8080"
)

// Handler handler is an interface that objects can implement to be registered to serve as middleware
// in the Negroni middleware stack.
// ServeHTTP should yield to the next middleware in the chain by invoking the next http.HandlerFunc
// passed in.
//
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.
type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFun