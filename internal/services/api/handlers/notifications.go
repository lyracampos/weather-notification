package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"weather-notification/internal/domain/usecases"

	"go.uber.org/zap"
)

type notificationHandler struct {
	log                       *zap.SugaredLogger
	queueNotificationsUseCase usecases.EnqueueNotificationsUseCase
}

func NewNotificationHandler(log *zap.SugaredLogger, queueNotificationsUseCase usecases.EnqueueNotificationsUseCase) *notificationHandler {
	return &notificationHandler{
		log:                       log,
		queueNotificationsUseCase: queueNotificationsUseCase,
	}
}

func (h *notificationHandler) Notify(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("notificationHandler.Notify - started")

	ctx := r.Context()
	rw.Header().Set("Content-type", "application/json")

	requestBody, _ := io.ReadAll(r.Body)
	var requestedInput usecases.QueueNotificationsInput
	if err := json.Unmarshal(requestBody, &requestedInput); err != nil {
		h.handlerErrors(rw, err)

		return
	}

	err := h.queueNotificationsUseCase.Execute(ctx, requestedInput)
	if err != nil {
		h.handlerErrors(rw, err)

		return
	}

	h.log.Info("notificationHandler.Notify - finished")

	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(requestedInput); err != nil {
		h.log.Error(fmt.Errorf("notificationHandler.Notify - encode failed: %w", err))
	}

}

// nolint: errcheck
func (h *notificationHandler) handlerErrors(rw http.ResponseWriter, err error) {
	h.log.Error(err.Error())

	switch {
	case strings.Contains(err.Error(), "Error:Field"):
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	default:
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
}
