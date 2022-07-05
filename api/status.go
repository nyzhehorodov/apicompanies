package api

import (
	"encoding/json"
	"net/http"

	"github.com/nyzhehorodov/apicompanies/pkg/version"
)

// StatusResponse holds application info and availability status.
type StatusResponse struct {
	App       string `json:"app"`
	Health    string `json:"health"`
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
}

func (a *API) statusHandler(w http.ResponseWriter, _ *http.Request) {
	jsResp, err := json.Marshal(
		StatusResponse{
			Health:    "ok",
			App:       version.AppName(),
			Version:   version.AppVersion(),
			GitCommit: version.Commit(),
			BuildDate: version.Date(),
		},
	)
	if err != nil {
		a.logger.Error(err, "serialize status response")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(jsResp); err != nil {
		a.logger.Error(err, "got an error processing response")
	}
}
