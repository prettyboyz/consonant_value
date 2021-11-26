package middleware

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type subrouter struct {
	c *web.C
	h http.Handler
}

func (s subrouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.c.URLParams != nil {
		path, ok := s.c.URLParams["*"]
		if !ok {
			path, ok = s.c.URLParams["_"]
		}
		if ok {
			oldpath := r.URL.Path
			oldmatch := web.GetMatch(*s.c)
			r.URL.Path = path
			if oldmatch.Handler != nil {
				delete(s.c.Env, web.MatchKey)
			}

			defer func() {
				r.URL.Path = oldpath

				if s.c.Env == nil {
					return
				}
				if oldmatch.Handler != nil {
					s.c.Env[web.MatchKey] = oldmatch
				} else {
					delete(s.c.Env, web.MatchKey)
				}
			}()
		}
	}
	s.h.ServeHTTP(w, r)
}

/*
SubRouter is a helper middleware that makes writing sub-routers easier.

If you register a sub-router under a key like "/admin/*", Goji's router will
automatically set c.URLParams["*"] to the unmatched path suffix. This middleware
will help you set the request URL's Pa