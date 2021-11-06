/*
Package cors is net/http handler to handle CORS related requests
as defined by http://www.w3.org/TR/cors/

You can configure it by passing an option struct to cors.New:

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"foo.com"},
        AllowedMethods: []string{"GET", "POST", "DELETE"},
        AllowCredentials: true,
    })

Then insert the handler in the chain:

    handler = c.Handler(handler)

See Options documentation for more options.

The resulting handler is a standard net/http handler.
*/
package cors

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/rs/xhandler"
	"golang.org/x/net/context"
)

// Options is a configuration container to setup the CORS middleware.
type Options struct {
	// AllowedOrigins is a list of origins a cross-domain reques