package websocket

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"nhooyr.io/websocket"
)

// type Client struct {
// 	Nickname string
// 	conn     *websocket.Conn
// 	ctx      context.Context
// }

var (
	clients2 map[*websocket.Conn]bool = make(map[*websocket.Conn]bool)
)

type wsServerHandler struct {
	log *zap.SugaredLogger
}

func NewWsServerHandler(log *zap.SugaredLogger) *wsServerHandler {
	return &wsServerHandler{
		log: log,
	}
}

func (h *wsServerHandler) ServerHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})

	if err != nil {
		h.log.Errorf("failed to start websocket server: %v", err)
	}

	clients2[conn] = true
	for {
		_, data, err := conn.Read(r.Context())
		if err != nil {
			h.log.Errorf("closing client connection: %v", err)
			delete(clients2, conn)
			break
		}
		h.log.Info(string(data))

		message := fmt.Sprintf("response from server %s", data)
		conn.Write(r.Context(), websocket.MessageText, []byte(message))
	}
}
