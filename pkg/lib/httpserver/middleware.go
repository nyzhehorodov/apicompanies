package httpserver

import (
	"net/http"
)

type middleware struct {
	f    func(http.HandlerFunc) http.HandlerFunc
	opts middlewareOptions
}

type middlewareOptions struct {
	prefix string
}

// MiddlewareOptions is a function on the options for a middleware.
type MiddlewareOption func(*middlewareOptions)

// PathPrefix is an Option to set the prefix
// for the routes to which this middleware should be applied
func PathPrefix(pref string) MiddlewareOption {
	return func(o *middlewareOptions) {
		o.prefix = pref
	}
}
