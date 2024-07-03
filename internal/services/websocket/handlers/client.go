package handlers

import (
	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

type Client struct {
	User string
	Conn *websocket.Conn
}

var (
	clients map[*Client]bool = make(map[*Client]bool)
)

type websocketHandler struct {
	log *zap.SugaredLogger
}

func NewWebSocketHandler(log *zap.SugaredLogger) *websocketHandler {
	return &websocketHandler{
		log: log,
	}
}
