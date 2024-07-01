package http

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"
	"weather-notification/configs"
	"weather-notification/internal/domain"
	"weather-notification/internal/domain/entities"
	"weather-notification/internal/domain/ports"
	"weather-notification/internal/gateways/http/models"

	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var _ ports.WeatherHTTPGateway = (*weatherAPI)(nil)

type weatherAPI struct {
	addressURL string
	timeout    int
	retries    int
	log        *zap.SugaredLogger
}

func NewWeatherAPI(log *zap.SugaredLogger, config *configs.Config) *weatherAPI {
	log.Info("weatherAPI - client started...")

	return &weatherAPI{
		addressURL: config.WeatherAPI.AddressURL,
		timeout:    config.WeatherAPI.Timeout,
		retries:    config.WeatherAPI.Retries,
		log:        log,
	}
}

func (w *weatherAPI) GetCity(ctx context.Context, city string, state string) (*entities.City, error) {
	url := fmt.Sprintf("%s/XML/listaCidades?city=%s", w.addressURL, city)
	response, err := w.doRequestWithRetry(url, http.MethodGet)
	if err != nil {
		return &entities.City{}, fmt.Errorf("failed to make HTTP request to weather API: %w", err)
	}

	decoder, err := w.readXMLKResponse(response)
	defer response.Body.Close()

	if err != nil {
		return &entities.City{}, fmt.Errorf("faield to read xml response from weather API: %v", err)
	}
	var cities models.Cities
	err = decoder.Decode(&cities)
	if err != nil {
		return &entities.City{}, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	for _, c := range cities.Cities {
		cityWithouAccents, err := removeAccents(c.Name)
		if err != nil {
			return &entities.City{}, fmt.Errorf("failed to validate city: %w", err)
		}
		if cityWithouAccents == city && c.State == state {
			return c.ToEntity(), nil
		}
	}

	return &entities.City{}, domain.ErrCityNotFound
}

func (w *weatherAPI) doRequestWithRetry(url string, httpMethod string) (*http.Response, error) {
	var (
		response *http.Response
		err      error
		attempts int = 1
	)
	httpClient := http.Client{Timeout: time.Duration(w.timeout) * time.Millisecond}

	url = strings.Replace(url, " ", "%20", -1)

	for w.retries > 0 {
		switch httpMethod {
		case http.MethodGet:
			response, err = httpClient.Get(url)
		default:
		}

		if err != nil {
			w.log.Errorf("attempt %d to make http request to the api failed %v:", attempts, err)
			w.retries -= 1
			attempts += 1
			time.Sleep(time.Duration(3) * time.Second)
		} else {
			break
		}
	}

	return response, err
}

func (w *weatherAPI) readXMLKResponse(response *http.Response) (*xml.Decoder, error) {
	if response.StatusCode != http.StatusOK {
		return &xml.Decoder{}, fmt.Errorf("status error: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return &xml.Decoder{}, fmt.Errorf("failed on reading response body: %v", err)

	}

	// dealing with ISO-8859-1
	reader := strings.NewReader(string(body))
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel

	return decoder, nil
}

func removeAccents(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	output, _, e := transform.String(t, s)
	if e != nil {
		panic(e)
	}
	return output, nil
}
