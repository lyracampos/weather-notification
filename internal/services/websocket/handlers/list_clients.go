package handlers

import (
	"encoding/json"
	"net/http"
)

// nolint: errcheck
func (h *websocketHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var res []string
	for c := range clients {
		res = append(res, c.User)
	}
	json.NewEncoder(w).Encode(res)
}
