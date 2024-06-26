package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"weather-notification/internal/domain"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/usecases"

	"go.uber.org/zap"
)

type userHandler struct {
	log             *zap.SugaredLogger
	registerUseCase usecases.RegisterUserUseCase
}

func NewUserHandler(log *zap.SugaredLogger, registerUseCase usecases.RegisterUserUseCase) *userHandler {
	return &userHandler{
		log:             log,
		registerUseCase: registerUseCase,
	}
}

// swagger:route POST /users  users RegisterUser
// Register new user in the application
// responses:
//
//	201: userRegisterResponse
//	400: notFoundResponse
//	501: internalServerErrorResponse
func (h *userHandler) Register(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("userHandler.Register - started")

	ctx := r.Context()
	rw.Header().Set("Content-type", "application/json")

	requestBody, _ := io.ReadAll(r.Body)
	var requestedUser entities.User
	if err := json.Unmarshal(requestBody, &requestedUser); err != nil {
		h.handlerErrors(rw, err)

		return
	}

	user, err := h.registerUseCase.Execute(ctx, &requestedUser)
	if err != nil {
		h.handlerErrors(rw, err)

		return
	}

	h.log.Info("userHandler.Register - finished successfully")

	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(user); err != nil {
		log.Printf("userHandler.Register - encode failed: %v", err)
	}
}

// nolint: errcheck
func (h *userHandler) handlerErrors(rw http.ResponseWriter, err error) {
	h.log.Error(err.Error())

	msgErr := err.Error()

	switch {
	case strings.Contains(msgErr, "Error:Field"):
		rw.WriteHeader(http.StatusBadRequest)
	case errors.Is(err, domain.ErrEmailIsAlreadyInUse):
		rw.WriteHeader(http.StatusConflict)
	default:
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Write([]byte(msgErr))
}
