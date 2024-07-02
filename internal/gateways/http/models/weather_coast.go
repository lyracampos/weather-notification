package models

import (
	"encoding/xml"
	"weather-notification/internal/domain/entities"
)

type WeatherCoast struct {
	XMLName   xml.Name         `xml:"cidade"`
	Morning   WeatherCoastData `xml:"manha"`
	Afternoon WeatherCoastData `xml:"tarde"`
	Evening   WeatherCoastData `xml:"noite"`
}

type WeatherCoastData struct {
	Day           string `xml:"dia"`
	SeaAgiation   string `xml:"agitacao"`
	WaveHeight    string `xml:"altura"`
	Direction     string `xml:"direcao"`
	WindSpeed     string `xml:"vento"`
	WindDirection string `xml:"vento_dir"`
}

func (w *WeatherCoast) ToEntity() *entities.WeatherCoast {
	return &entities.WeatherCoast{
		Morning: entities.WeatherCoastData{
			Day:           w.Morning.Day,
			SeaAgiation:   w.Morning.SeaAgiation,
			WaveHeight:    w.Morning.WaveHeight,
			Direction:     w.Morning.Direction,
			WindSpeed:     w.Morning.WindSpeed,
			WindDirection: w.Morning.WindDirection,
		},
	}
}
