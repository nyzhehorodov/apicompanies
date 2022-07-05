package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/nyzhehorodov/apicompanies/pkg/lib/httpserver"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/log"
)

// API represents apicompany app
type API struct {
	server *httpserver.Server
	logger log.Interface
}

// New returns new instance of API
func New(server *httpserver.Server, logger log.Interface) (*API, error) {
	return &API{
		server: server,
		logger: logger,
	}, nil
}

// Config holds parameters for API
type Config struct {
	Port    int
	TLSKey  string
	TLSCert string
}

// Run starts application and blocks the caller.
func (a *API) Run(ctx context.Context, c Config) {
	listenAddr := fmt.Sprintf(":%d", c.Port)
	a.logger.Info("server listen", "addr", listenAddr)

	a.initHandlers()

	go func() {
		<-ctx.Done()
		if err := a.server.Shutdown(context.Background()); err != nil {
			a.logger.Error(err, "server shutdown")
		}
	}()

	var err error
	if c.TLSCert != "" && c.TLSKey != "" {
		err = a.server.ListenTLS(listenAddr, c.TLSCert, c.TLSKey)
	} else {
		err = a.server.Listen(listenAddr)
	}

	if !errors.Is(err, http.ErrServerClosed) {
		a.logger.Error(err, "server listen")
	}
}

func (a *API) initHandlers() {
	a.server.HandleGET("/v1/status", a.statusHandler)
}
