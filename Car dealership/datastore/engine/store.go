package engine

import (
	"context"
	"database/sql"

	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/models"

	"github.com/google/uuid"
)

type Enginestore struct {
	db *sql.DB
}

func New(db *sql.DB) Enginestore {
	return Enginestore{db: db}
}

// EngineGetByID store layer function to get engine details
func (s Enginestore) EngineGetByID(ctx context.Context, id string) (models.Engine, error) {
	var engine models.Engine

	err := s.db.QueryRowContext(ctx, "SELECT  *from Engine where id=?", id).
		Scan(&engine.EngineID, &engine.Displacement, &engine.NoOfCylinder, &engine.CarRange)
	if err != nil {
		return models.Engine{}, err
	}

	return engine, nil
}

// EngineCreate store layer function to create engine
func (s Enginestore) EngineCreate(ctx context.Context, engine *models.Engine) (models.Engine, error) {
	engine.EngineID = uuid.New()

	_, err := s.db.ExecContext(ctx, "INSERT INTO Engine (id,displacement,cylinders,`range`) VALUES(?,?,?,?)",
		engine.EngineID.String(), engine.Displacement, engine.NoOfCylinder, engine.CarRange)
	if err != nil {
		return models.Engine{}, err
	}

	return *engine, nil
}

// EngineUpdate store layer function to update engine details
func (s Enginestore) EngineUpdate(ctx context.Context, id string, engine models.Engine) (models.Engine, error) {
	_, err := s.db.ExecContext(ctx,
		"UPDATE Engine SET displacement=?,cylinders=?,`range`=? WHERE Id=?",
		engine.Displacement, engine.NoOfCylinder, engine.CarRange, id)
	if err != nil {
		return models.Engine{}, err
	}

	engine.EngineID = uuid.MustParse(id)

	return engine, nil
}

// EngineDelete to delete engine record
func (s Enginestore) EngineDelete(ctx context.Context, id string) (models.Engine, error) {
	_, err := s.db.ExecContext(ctx, "delete from Engine where id=?", id)
	if err != nil {
		return models.Engine{}, err
	}

	return models.Engine{}, nil
}
