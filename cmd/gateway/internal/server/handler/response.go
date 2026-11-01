package handler

import (
	"encoding/json"
	"net/http"
)

type response struct {
	Status int          `json:"status"`
	Body   responseBody `json:"body"`
}

type responseBody struct {
	Info map[string]any `json:"info,omitempty"`
}

func (cr *CustomRouter) mustSendResponseJSON(w http.ResponseWriter, status int, bodyInfo map[string]any) {
	const fnName = "send response json"

	resp := response{
		Status: status,
	}

	if len(bodyInfo) != 0 {
		resp.Body.Info = bodyInfo
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		cr.log.Error("Can't send response", cr.log.Attr("error", err), cr.log.Attr("func", fnName))
	}
}
