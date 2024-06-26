package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"weather-notification/internal/domain"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/usecases"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

var (
	ErrUserDoesNotExist  = errors.New("no user was found for this email. Please check the email and try again")
	ErrEmailAlreadyInUse = errors.New("email in use by another user")
)

type userHandler struct {
	log                *zap.SugaredLogger
	registerUseCase    usecases.RegisterUserUseCase
	unsubscribeUseCase usecases.UnsubscribeUserUseCase
}

func NewUserHandler(log *zap.SugaredLogger, registerUseCase usecases.RegisterUserUseCase, unsubscribeUseCase usecases.UnsubscribeUserUseCase) *userHandler {
	return &userHandler{
		log:                log,
		registerUseCase:    registerUseCase,
		unsubscribeUseCase: unsubscribeUseCase,
	}
}

// swagger:route POST /users  users RegisterUser
// Register new user in the application
// responses:
//
//	201: userRegisterResponse
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

	h.log.Info("userHandler.Register - finished")

	rw.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(rw).Encode(user); err != nil {
		h.log.Error(fmt.Errorf("userHandler.Register - encode failed: %w", err))
	}
}

// swagger:route PUT /users  users Unsubscribe
// Unsubscribe user from list of notifications
// responses:
//
//	200: userUnsubscribeResponse
//	404: notFoundResponse
//	501: internalServerErrorResponse
func (h *userHandler) Unsubscribe(rw http.ResponseWriter, r *http.Request) {
	h.log.Info("userHandler.Unsubscribe - started")

	ctx := r.Context()
	rw.Header().Set("Content-type", "application/json")

	vars := mux.Vars(r)
	email := vars["email"]

	user, err := h.unsubscribeUseCase.Execute(ctx, email)
	if err != nil {
		h.handlerErrors(rw, err)

		return
	}

	h.log.Info("userHandler.Unsubscribe - finished")

	rw.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(rw).Encode(user); err != nil {
		h.log.Error(fmt.Errorf("userHandler.Unsubscribe - encode failed: %w", err))
	}
}

// nolint: errcheck
func (h *userHandler) handlerErrors(rw http.ResponseWriter, err error) {
	h.log.Error(err.Error())

	switch {
	case strings.Contains(err.Error(), "Error:Field"):
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
	case errors.Is(err, domain.ErrUserNotFound):
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(ErrUserDoesNotExist.Error()))
	case errors.Is(err, domain.ErrEmailIsAlreadyInUse):
		rw.WriteHeader(http.StatusConflict)
		rw.Write([]byte(ErrEmailAlreadyInUse.Error()))
	default:
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
}
