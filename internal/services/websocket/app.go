package websocket

import (
	"encoding/json"
	"io"
	"net/http"

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

func (h *websocketHandler) Clients(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var res []string
	for c := range clients {
		res = append(res, c.User)
	}
	json.NewEncoder(w).Encode(res)
}

type Notification struct {
	Message string `json:"message"`
}

func (h *websocketHandler) NotifyUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	user := r.URL.Query().Get("user")

	requestBody, _ := io.ReadAll(r.Body)
	var notification Notification
	if err := json.Unmarshal(requestBody, &notification); err != nil {
		h.log.Errorf("failed to unmarshal JSON message to struct: %v", err)
		return
	}

	for client, c := range clients {
		if client.User == user && c {
			client.Conn.Write(r.Context(), websocket.MessageText, []byte(notification.Message))
			break
		}
	}
	res, _ := json.Marshal(notification.Message)
	w.Write(res)
}
