package models

import (
	"encoding/xml"
	"weather-notification/internal/domain/entities"
)

type Cities struct {
	XMLName xml.Name `xml:"cidades"`
	Cities  []City   `xml:"cidade"`
}

type City struct {
	ID    int    `xml:"id"`
	Name  string `xml:"nome"`
	State string `xml:"uf"`
}

func (c City) ToEntity() *entities.City {
	return &entities.City{
		ID:    c.ID,
		Name:  c.Name,
		State: c.State,
	}
}
