package api

import (
	"encoding/json"
	"net/http"
)

func decodeRequest(r *http.Request, req interface{}) error {
	return json.NewDecoder(r.Body).Decode(req)
}

func encodeResponse(w http.ResponseWriter, resp interface{}) error {
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(resp)
}
