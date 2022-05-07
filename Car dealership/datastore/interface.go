package datastore

import (
	"context"

	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/models"
)

type Car interface {
	GetCarByID(ctx context.Context, id string) (models.Car, error)
	GetCarsByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error)
	CreateCar(ctx context.Context, car *models.Car) (models.Car, error)
	DeleteCar(ctx context.Context, id string) (models.Car, error)
	UpdateCar(ctx context.Context, id string, car models.Car) (models.Car, error)
}

type Engine interface {
	EngineGetByID(ctx context.Context, id string) (models.Engine, error)
	EngineCreate(ctx context.Context, engine *models.Engine) (models.Engine, error)
	EngineDelete(ctx context.Context, id string) (models.Engine, error)
	EngineUpdate(ctx context.Context, id string, engine models.Engine) (models.Engine, error)
}
