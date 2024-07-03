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

	"github.com/cenkalti/backoff"
	"go.uber.org/zap"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var _ ports.WeatherHTTPGateway = (*weatherAPI)(nil)

const (
	DefaultInitialInterval     = 500 * time.Millisecond
	DefaultRandomizationFactor = 0.5
	DefaultMultiplier          = 1.5
	DefaultMaxInterval         = 60 * time.Second
	DefaultMaxElapsedTime      = 15 * time.Minute
)

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
	var (
		response *http.Response
		err      error
	)

	url := fmt.Sprintf("%s/XML/listaCidades?city=%s", w.addressURL, city)
	getCity := func() error {
		w.log.Infof("retry to get city")
		response, err = w.doRequest(url, http.MethodGet)
		if err != nil {
			return fmt.Errorf("failed to get city from weather API %v", err)
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second

	err = backoff.Retry(getCity, expBackoff)
	if err != nil {
		return &entities.City{}, fmt.Errorf("failed to connect to get city from weather API after retrying: %v", err)
	}

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

func (w *weatherAPI) GetWeather(ctx context.Context, cityID int) (*[]entities.Weather, error) {
	var (
		response *http.Response
		err      error
	)

	url := fmt.Sprintf("%s/XML/cidade/%d/previsao.xml", w.addressURL, cityID)
	getWeather := func() error {
		w.log.Infof("retry to get city")
		response, err = w.doRequest(url, http.MethodGet)
		if err != nil {
			return fmt.Errorf("failed to get weather from weather API %v", err)
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second

	err = backoff.Retry(getWeather, expBackoff)
	if err != nil {
		return &[]entities.Weather{}, fmt.Errorf("failed to connect to get wather from weather API after retrying: %v", err)
	}

	if err != nil {
		return &[]entities.Weather{}, fmt.Errorf("failed to make HTTP request to weather API: %v", err)
	}

	decoder, err := w.readXMLKResponse(response)
	defer response.Body.Close()

	if err != nil {
		return &[]entities.Weather{}, fmt.Errorf("faield to read xml response from weather API: %v", err)
	}

	var weatherList models.WeatherList
	err = decoder.Decode(&weatherList)
	if err != nil {
		return &[]entities.Weather{}, fmt.Errorf("failed to unmarshal XML: %v", err)
	}

	return weatherList.ToEntity(), nil
}

func (w *weatherAPI) GetWeatherCoast(ctx context.Context, cityID int) (*entities.WeatherCoast, error) {
	var (
		response *http.Response
		err      error
	)

	url := fmt.Sprintf("%s/XML/cidade/%d/dia/0/ondas.xml", w.addressURL, cityID)

	getWeatherCoast := func() error {
		w.log.Infof("retry to get city")
		response, err = w.doRequest(url, http.MethodGet)
		if err != nil {
			return fmt.Errorf("failed to get weather coast from weather API %v", err)
		}

		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second

	err = backoff.Retry(getWeatherCoast, expBackoff)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to get wather coast from weather API after retrying: %v", err)
	}
	// response, err := w.doRequestWithRetry2(url, http.MethodGet)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request to weather API: %v", err)
	}

	decoder, err := w.readXMLKResponse(response)
	defer response.Body.Close()

	if err != nil {
		return nil, fmt.Errorf("faield to read xml response from weather API: %v", err)
	}

	var weatherList models.WeatherCoast
	err = decoder.Decode(&weatherList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %v", err)
	}

	if weatherList.Morning.SeaAgiation == "undefined" {
		return nil, nil
	}

	return weatherList.ToEntity(), nil
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

func (w *weatherAPI) doRequest(url string, httpMethod string) (*http.Response, error) {
	var (
		response *http.Response
		err      error
	)

	httpClient := http.Client{Timeout: time.Duration(w.timeout) * time.Millisecond}
	url = strings.Replace(url, " ", "%20", -1)

	switch httpMethod {
	case http.MethodGet:
		response, err = httpClient.Get(url)
	default:
	}

	return response, err
}
