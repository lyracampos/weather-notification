// nolint:unused
// Package Weather Notification API
//
// Documentation for Weather Notification API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package docs

import (
	"weather-notification/internal/domain/entities"
)

// Internal server error message returned as a string
// swagger:response internalServerErrorResponse
type internalServerErrorResponseWrapper struct {
	// error description
	// in: body
	Body MessageError
}

// Not found message error returned as string
// swagger:response notFoundResponse
type errorNotFoundResponseWrapper struct {
	// error description
	// in: body
	Body MessageError
}

// BadRequest message error returned as string
// swagger:response badRequestResponse
type badRequestResponseWrapper struct {
	// error description
	// in: body
	Body MessageError
}

type MessageError struct {
	Message string `json:"message"`
}

// Data structure representing user registered
// swagger:response userRegisterResponse
type userRegisterResponseWrapper struct {
	// in: body
	Body entities.User
}

// swagger:parameters Register
type userRegisterCommandWrapper struct {
	// Payload to register new user in application
	// in: body
	// required: true
	Body entities.User
}

// Data structure representing user unsubscribing
// swagger:response userUnsubscribeResponse
type userUnsubscribeResponseWrapper struct {
	// in: body
	Body entities.User
}
