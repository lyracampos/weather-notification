package handlers

import (
	"net/http"

	"go.uber.org/zap"
)

type healthHandler struct {
	log *zap.SugaredLogger
}

func NewHealthHandler(log *zap.SugaredLogger) *healthHandler {
	return &healthHandler{
		log: log,
	}
}

// nolint: errcheck
func (h *healthHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.log.Info("HandlerHealth - checking API status")
	w.Write([]byte("user API is running"))
}
