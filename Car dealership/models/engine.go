package models

import "github.com/google/uuid"

type Engine struct {
	EngineID     uuid.UUID `json:"id"`
	Displacement int64     `json:"Displacement"`
	NoOfCylinder int64     `json:"NoOfCylinder"`
	CarRange     int64     `json:"Range"`
}
