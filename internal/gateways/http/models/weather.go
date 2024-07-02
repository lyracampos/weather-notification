package models

import (
	"encoding/xml"
	"weather-notification/internal/domain/entities"
)

type WeatherList struct {
	XMLName     xml.Name  `xml:"cidade"`
	WeatherList []Weather `xml:"previsao"`
}

type Weather struct {
	Day       string  `xml:"dia"`
	Condition string  `xml:"tempo"`
	Max       int     `xml:"maxima"`
	Min       int     `xml:"minima"`
	IUV       float32 `xml:"iuv"`
}

func (w *WeatherList) ToEntity() *[]entities.Weather {
	result := make([]entities.Weather, 0, len(w.WeatherList))
	for _, item := range w.WeatherList {
		weather := entities.Weather{
			Day:       item.Day,
			Condition: item.Condition,
			Max:       item.Max,
			Min:       item.Min,
			IUV:       item.IUV,
		}
		result = append(result, weather)
	}
	return &result
}
