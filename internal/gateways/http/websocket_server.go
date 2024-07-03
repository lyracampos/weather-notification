package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"weather-notification/configs"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"

	"go.uber.org/zap"
)

var _ ports.WebNotificationHTTPGateway = (*webSocketServer)(nil)

type webSocketServer struct {
	addressURL string
	timeout    int
	log        *zap.SugaredLogger
}

func NewWebSocketServer(log *zap.SugaredLogger, config *configs.Config) *webSocketServer {
	log.Info("websocket server - started...")

	return &webSocketServer{
		addressURL: fmt.Sprintf("http://%s:%d", config.WebSocketClient.Host, config.WebSocketClient.Port),
		timeout:    config.WeatherAPI.Timeout,
		log:        log,
	}
}

type notification struct {
	Weather      *[]entities.Weather
	WeatherCoast *entities.WeatherCoast
}

func (w *webSocketServer) SendNotification(ctx context.Context, user *entities.User, weather *[]entities.Weather, weatherCoast *entities.WeatherCoast) error {
	url := fmt.Sprintf("%s/ws/notify?user=%s", w.addressURL, user.Email)

	payload := notification{
		Weather:      weather,
		WeatherCoast: weatherCoast,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to convert struct to json: %w", err)
	}

	reader := bytes.NewReader(body)

	resp, err := http.Post(url, "application/json", reader)
	if err != nil {
		return fmt.Errorf("failed to make post request to websocket server")
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			w.log.Errorf("failed to close response from websocket server: %v", err)
		}
	}()

	if resp.StatusCode >= 400 && resp.StatusCode <= 500 {
		w.log.Error("failed to send notification - status Code: ", resp.StatusCode)
	}

	return nil
}
