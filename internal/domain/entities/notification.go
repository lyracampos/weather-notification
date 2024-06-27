package entities

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Type uint8

const (
	NotificationTypeUnspecified = "unspecified"
	NotificationTypeWebSocket   = "websocket"
	NotificationTypePush        = "push"
	NotificationTypeEmail       = "email"
	NotificationTypeSMS         = "sms"

	NotificationStatusUnspecified = "unespecified"
	NotificationStatusQueued      = "queued"
	NotificationStatusSent        = "sent"
)

// swagger:model
type Notification struct {
	// notification identification ID
	//
	ID int64
	// the email of the user who will be notified
	//
	// required: true
	EmailTo string `validate:"required,email" json:"email_to"`
	// the notification sending type (web, push, email or sms)
	//
	// required: true
	Type string `validate:"required,oneof=web push email sms"`
	// notification sending status (queued or sent)
	//
	// required: true
	Status string `validate:"required,oneof=queued sent"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewNotification(emailTo, notificationType, status string) *Notification {
	return &Notification{
		EmailTo: emailTo,
		Type:    notificationType,
		Status:  status,
	}
}

func (u *Notification) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}
