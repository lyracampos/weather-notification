package events

type WebsocketNotificationEvent struct {
	UserEmail string
}

type EmailNotificationEvent struct {
	UserEmail string
}

type SMSNotificationEvent struct {
	UserEmail string
}

type PushNotificationEvent struct {
	UserEmail string
}
