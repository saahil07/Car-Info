package car

import (
	"context"
	"database/sql"

	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/models"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return Store{db: db}
}

// GetCarByID store layer function to get car details when car id is provided
func (s Store) GetCarByID(ctx context.Context, id string) (models.Car, error) {
	var c models.Car

	err := s.db.QueryRowContext(ctx, "SELECT * FROM Car WHERE ID=?;", id).
		Scan(&c.ID, &c.Engine.EngineID, &c.Name, &c.Year, &c.Brand, &c.FuelType)
	if err != nil {
		return models.Car{}, err
	}

	return c, nil
}

// GetCarsByBrand store layer function to get all car records of brand name given
func (s Store) GetCarsByBrand(ctx context.Context, brand string, isEngine bool) ([]models.Car, error) {
	rows, err := s.db.QueryContext(ctx, "select * from Car where brand=?;", brand)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var car []models.Car

	for rows.Next() {
		var c models.Car

		err = rows.Scan(&c.ID, &c.Engine.EngineID, &c.Name, &c.Year, &c.Brand, &c.FuelType)
		if err != nil {
			return nil, err
		}

		car = append(car, c)
	}

	err = rows.Err()
	if err != nil {
		return []models.Car{}, err
	}

	return car, nil
}

// CreateCar store layer function to create car record
func (s Store) CreateCar(ctx context.Context, car *models.Car) (models.Car, error) {
	_, err := s.db.ExecContext(ctx, "INSERT INTO Car (id,engine_id,name,year,brand,fuel_type) VALUES(?,?,?,?,?,?)",
		car.ID.String(), car.Engine.EngineID, car.Name, car.Year, car.Brand, car.FuelType)
	if err != nil {
		return models.Car{}, err
	}

	return *car, nil
}

// UpdateCar store layer function to update car record
func (s Store) UpdateCar(ctx context.Context, id string, car models.Car) (models.Car, error) {
	_, err := s.db.ExecContext(ctx, "UPDATE Car SET name=?,year=?,brand=?,fuel_type=? WHERE id=?",
		car.Name, car.Year, car.Brand, car.FuelType, id)
	if err != nil {
		return models.Car{}, err
	}

	car.ID = uuid.MustParse(id)

	return car, nil
}

// DeleteCar store layer function to delete car record
func (s Store) DeleteCar(ctx context.Context, id string) (models.Car, error) {
	_, err := s.db.ExecContext(ctx, "DELETE FROM Car WHERE ID=?", id)
	if err != nil {
		return models.Car{}, err
	}

	return models.Car{}, nil
}
