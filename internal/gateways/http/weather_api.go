package http

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"weather-notification/configs"
	"weather-notification/internal/domain"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"
	"weather-notification/internal/gateways/http/models"

	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
)

var _ ports.WeatherHTTPGateway = (*weatherAPI)(nil)

type weatherAPI struct {
	addressURL string
}

func NewWeatherAPI(log *zap.SugaredLogger, config *configs.Config) *weatherAPI {
	log.Info("weatherAPI - client started")

	return &weatherAPI{
		addressURL: config.WeatherAPI.AddressURL,
	}
}

func (w *weatherAPI) GetCity(ctx context.Context, city string, state string) (*entities.City, error) {
	url := fmt.Sprintf("%s/XML/listaCidades?city=%s", w.addressURL, city)
	resp, err := http.Get(url)
	if err != nil {
		return &entities.City{}, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &entities.City{}, fmt.Errorf("status error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &entities.City{}, fmt.Errorf("failed on reading response body: %w", err)

	}

	// dealing with ISO-8859-1
	reader := strings.NewReader(string(body))
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	var cities models.Cities
	err = decoder.Decode(&cities)
	if err != nil {
		return &entities.City{}, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	for _, c := range cities.Cities {
		if c.Name == city && c.State == state {
			return c.ToEntity(), nil
		}
	}

	return &entities.City{}, domain.ErrCityNotFound
}
