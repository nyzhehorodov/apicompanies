package api

import (
	"net/http"

	"github.com/nyzhehorodov/apicompanies/api/v1"
	"github.com/nyzhehorodov/apicompanies/pkg/domain/company"
	"github.com/nyzhehorodov/apicompanies/pkg/lib/ctxparam"
)

func (a *API) CompanyAddHandler(w http.ResponseWriter, r *http.Request) {
	req := &v1.CompanyRequest{}
	if err := decodeRequest(r, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comp := company.Company{
		Code:    req.Code,
		Name:    req.Name,
		Country: req.Country,
		Website: req.Website,
		Phone:   req.Phone,
	}

	err := a.CompanyService.Add(comp)
	if err != nil {
		a.Logger.Error(err, "handler add company")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) CompanyListHandler(w http.ResponseWriter, r *http.Request) {
	//	TODO
	//	not implemented
}

func (a *API) CompanyUpdateHandler(w http.ResponseWriter, r *http.Request) {
	req := &v1.CompanyRequest{}
	if err := decodeRequest(r, req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	comp := company.Company{
		Code:    req.Code,
		Name:    req.Name,
		Country: req.Country,
		Website: req.Website,
		Phone:   req.Phone,
	}

	err := a.CompanyService.Update(&comp)
	if err != nil {
		a.Logger.Error(err, "handler update company")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *API) CompanyDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := ctxparam.Int(r.Context(), "id")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = a.CompanyService.Delete(id)
	if err != nil {
		a.Logger.Error(err, "handler delete company")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
