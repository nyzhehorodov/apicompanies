package api

import (
	"github.com/nyzhehorodov/apicompanies/pkg/app/company"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/httpserver"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/log"
)

// API represents apicompany app
type API struct {
	Server         *httpserver.Server
	Logger         log.Interface
	CompanyService company.Service
}

func (a *API) Init() {
	a.Server.HandleGET("/v1/status", a.StatusHandler)
	a.Server.HandlePOST("/v1/status", a.CompanyAddHandler)
	a.Server.HandleGET("/v1/company", a.CompanyListHandler)
	a.Server.HandlePUT("/v1/status", a.CompanyUpdateHandler)
	a.Server.HandleDELETE("/v1/status", a.CompanyDeleteHandler)
}
