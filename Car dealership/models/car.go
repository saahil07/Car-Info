package models

import "github.com/google/uuid"

type Car struct {
	ID       uuid.UUID `json:"ID"`
	Name     string    `json:"Name"`
	Year     int       `json:"Year"`
	Brand    string    `json:"Brand"`
	FuelType string    `json:"FuelType"`
	Engine   Engine    `json:"Engine"`
}
