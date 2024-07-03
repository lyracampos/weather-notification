package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"weather-notification/internal/domain/entities"

	"nhooyr.io/websocket"
)

type notification struct {
	Weather      *[]entities.Weather
	WeatherCoast *entities.WeatherCoast
}

func (h *websocketHandler) NotifyUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	user := r.URL.Query().Get("user")

	requestBody, _ := io.ReadAll(r.Body)
	var notificationByte notification
	if err := json.Unmarshal(requestBody, &notificationByte); err != nil {
		h.log.Errorf("failed to unmarshal JSON message to struct: %v", err)
		return
	}

	for client, c := range clients {
		if client.User == user && c {
			client.Conn.Write(r.Context(), websocket.MessageText, requestBody)
			break
		}
	}
	res, _ := json.Marshal(requestBody)
	w.Write(res)
}
