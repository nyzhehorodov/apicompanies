// Package httpserver encapsulates handy methods for running a REST-enabled API endpoint
package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/nyzhehorodov/apicompanies/pkg/lib/ctxparam"
)

// Server main object
type Server struct {
	router     *httprouter.Router
	middleware []middleware

	httpserver *http.Server

	done chan struct{}
}

// New returns a new server instance.
func New() *Server {
	router := httprouter.New()

	router.HandleMethodNotAllowed = false

	srv := &Server{
		router: router,
		done:   make(chan struct{}),
	}

	srv.SetNotFoundHandler(notImplementedHandler)

	return srv
}

// Close closes server with interrupting active connections.
func (srv *Server) Close() error {
	if srv.isClosed() {
		return nil
	}

	close(srv.done)

	err := srv.httpserver.Close()
	if err != nil {
		return fmt.Errorf("close http server: %s", err)
	}

	return nil
}

// Shutdown gracefully shuts down the server without interrupting any active connections.
func (srv *Server) Shutdown(ctx context.Context) error {
	if srv.isClosed() {
		return nil
	}

	close(srv.done)

	err := srv.httpserver.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("shutdown http server: %s", err)
	}
	return nil
}

func (srv *Server) isClosed() bool {
	select {
	case <-srv.done:
		return true
	default:
		return false
	}
}

// SetNotFoundHandler sets a custom handler for unregistered paths.
// The default one returns http 501 error (Not Implemented)
func (srv *Server) SetNotFoundHandler(handler http.HandlerFunc) {
	srv.router.NotFound = handler
}

// AddMiddleware adds middleware functions to the server.
// They will be executed for each call in the registration order.
func (srv *Server) AddMiddleware(handler func(http.HandlerFunc) http.HandlerFunc, opts ...MiddlewareOption) {
	m := middleware{f: handler}
	for _, opt := range opts {
		opt(&m.opts)
	}

	srv.middleware = append(srv.middleware, m)
}

// Listen starts serving at specified address and port.
// Always returns not nil error.
func (srv *Server) Listen(addr string) error {
	return srv.listen(addr, func(server *http.Server) error { return server.ListenAndServe() })
}

// ListenTLS starts serving at specified address and port.
// Always returns not nil error.
func (srv *Server) ListenTLS(addr, certFile, keyFile string) error {
	return srv.listen(addr, func(server *http.Server) error { return server.ListenAndServeTLS(certFile, keyFile) })
}

func (srv *Server) listen(addr string, listener func(server *http.Server) error) error {
	server := &http.Server{Addr: addr, Handler: srv.router}
	srv.httpserver = server
	return listener(server)
}

// HandleGET adds a new GET handler to the Server.
func (srv *Server) HandleGET(path string, handler http.HandlerFunc) {
	srv.handleFunc(path, handler, http.MethodGet)
}

// HandlePUT adds a new PUT handler to Server.
func (srv *Server) HandlePUT(path string, handler http.HandlerFunc) {
	srv.handleFunc(path, handler, http.MethodPut)
}

// HandlePOST adds a new POST handler to Server.
func (srv *Server) HandlePOST(path string, handler http.HandlerFunc) {
	srv.handleFunc(path, handler, http.MethodPost)
}

// HandleDELETE adds a new DELETE handler to Server.
func (srv *Server) HandleDELETE(path string, handler http.HandlerFunc) {
	srv.handleFunc(path, handler, http.MethodDelete)
}

// HandleOPTIONS adds a new OPTIONS handler to Server.
func (srv *Server) HandleOPTIONS(path string, handler http.HandlerFunc) {
	srv.handleFunc(path, handler, http.MethodOptions)
}

// Handle adds a new handler to Server (GET, POST, PATCH, PUT, DELETE, HEAD)
func (srv *Server) Handle(path string, handler http.HandlerFunc) {
	srv.handleFunc(path, handler, http.MethodGet)
	srv.handleFunc(path, handler, http.MethodPost)
	srv.handleFunc(path, handler, http.MethodPatch)
	srv.handleFunc(path, handler, http.MethodPut)
	srv.handleFunc(path, handler, http.MethodDelete)
	srv.handleFunc(path, handler, http.MethodHead)
}

func (srv *Server) handleFunc(path string, handler http.HandlerFunc, method string) {
	// apply middleware
	for i := len(srv.middleware) - 1; i >= 0; i-- {
		m := srv.middleware[i]
		if m.opts.prefix != "" && !strings.HasPrefix(path, m.opts.prefix) {
			continue
		}
		handler = m.f(handler)
	}

	h := paramsMiddleware(handler)

	srv.router.Handle(method, path, h)
}

func paramsMiddleware(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := r.Context()
		for _, p := range ps {
			ctx = ctxparam.WithValue(ctx, p.Key, p.Value)
		}

		paramsRequest := r.WithContext(ctx)
		next(w, paramsRequest)
	}
}

func notImplementedHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServeFiles adds a new handler that serves static files
func (srv *Server) ServeFiles(urlprefix, rootpath string) {
	urlprefix = strings.TrimRight(urlprefix, "/")
	srv.router.ServeFiles(urlprefix+"/*filepath", http.Dir(rootpath))
}

// ServeFile adds a new handler to serve single file.
func (srv *Server) ServeFile(urlpath, filepath string) {
	srv.router.HandlerFunc(http.MethodGet, urlpath, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath)
	})
}
