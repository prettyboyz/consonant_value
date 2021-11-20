
// Command example is a sample application built with Goji. Its goal is to give
// you a taste for what Goji looks like in the real world by artificially using
// all of its features.
//
// In particular, this is a complete working site for gritter.com, a site where
// users can post 140-character "greets". Any resemblance to real websites,
// alive or dead, is purely coincidental.
package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/goji/param"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

// Note: the code below cuts a lot of corners to make the example app simple.

func main() {
	// Add routes to the global handler
	goji.Get("/", Root)
	// Fully backwards compatible with net/http's Handlers
	goji.Get("/greets", http.RedirectHandler("/", 301))
	// Use your favorite HTTP verbs
	goji.Post("/greets", NewGreet)
	// Use Sinatra-style patterns in your URLs
	goji.Get("/users/:name", GetUser)
	// Goji also supports regular expressions with named capture groups.
	goji.Get(regexp.MustCompile(`^/greets/(?P<id>\d+)$`), GetGreet)

	// Middleware can be used to inject behavior into your app. The
	// middleware for this application are defined in middleware.go, but you
	// can put them wherever you like.
	goji.Use(PlainText)

	// If the patterns ends with "/*", the path is treated as a prefix, and
	// can be used to implement sub-routes.
	admin := web.New()
	goji.Handle("/admin/*", admin)

	// The standard SubRouter middleware helps make writing sub-routers
	// easy. Ordinarily, Goji does not manipulate the request's URL.Path,
	// meaning you'd have to repeat "/admin/" in each of the following
	// routes. This middleware allows you to cut down on the repetition by
	// eliminating the shared, already-matched prefix.
	admin.Use(middleware.SubRouter)
	// You can also easily attach extra middleware to sub-routers that are
	// not present on the parent router. This one, for instance, presents a
	// password prompt to users of the admin endpoints.
	admin.Use(SuperSecure)

	admin.Get("/", AdminRoot)
	admin.Get("/finances", AdminFinances)

	// Goji's routing, like Sinatra's, is exact: no effort is made to
	// normalize trailing slashes.
	goji.Get("/admin", http.RedirectHandler("/admin/", 301))

	// Use a custom 404 handler
	goji.NotFound(NotFound)

	// Sometimes requests take a long time.
	goji.Get("/waitforit", WaitForIt)

	// Call Serve() at the bottom of your main() function, and it'll take
	// care of everything else for you, including binding to a socket (with
	// automatic support for systemd and Einhorn) and supporting graceful
	// shutdown on SIGINT. Serve() is appropriate for both development and
	// production.
	goji.Serve()
}

// Root route (GET "/"). Print a list of greets.
func Root(w http.ResponseWriter, r *http.Request) {
	// In the real world you'd probably use a template or something.
	io.WriteString(w, "Gritter\n======\n\n")
	for i := len(Greets) - 1; i >= 0; i-- {
		Greets[i].Write(w)
	}
}

// NewGreet creates a new greet (POST "/greets"). Creates a greet and redirects
// you to the created greet.
//
// To post a new greet, try this at a shell:
// $ now=$(date +'%Y-%m-%dT%H:%M:%SZ')
// $ curl -i -d "user=carl&message=Hello+World&time=$now" localhost:8000/greets
func NewGreet(w http.ResponseWriter, r *http.Request) {
	var greet Greet

	// Parse the POST body into the Greet struct. The format is the same as
	// is emitted by (e.g.) jQuery.param.
	r.ParseForm()
	err := param.Parse(r.Form, &greet)

	if err != nil || len(greet.Message) > 140 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// We make no effort to prevent races against other insertions.
	Greets = append(Greets, greet)
	url := fmt.Sprintf("/greets/%d", len(Greets)-1)
	http.Redirect(w, r, url, http.StatusCreated)
}

// GetUser finds a given user and her greets (GET "/user/:name")
func GetUser(c web.C, w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Gritter\n======\n\n")
	handle := c.URLParams["name"]
	user, ok := Users[handle]
	if !ok {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	user.Write(w, handle)