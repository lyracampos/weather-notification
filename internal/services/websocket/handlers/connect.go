package handlers

import (
	"net/http"

	"nhooyr.io/websocket"
)

func (h *websocketHandler) Connect(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		h.log.Errorf("failed to start websocket server: %v", err)
	}

	client := &Client{User: user, Conn: conn}
	clients[client] = true

	h.log.Infof("the user %s has connected to websocket server", user)

	for {
		_, data, err := conn.Read(r.Context())
		if err != nil {
			h.log.Errorf("closing client connection: %v", err)
			delete(clients, client)
			break
		}

		h.log.Info(string(data))
	}
}
