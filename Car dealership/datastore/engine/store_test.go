package engine

import (
	"context"
	"errors"
	"testing"

	"github.com/zopsmart/GoLang-Interns-2022/tree/sahil-zs/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

// TestEnginestore_EngineGetByID function to test Enginegetbyid function
func TestEnginestore_EngineGetByID(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	dbcheck := New(db)

	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	engine := models.Engine{EngineID: id, Displacement: 1800, NoOfCylinder: 7, CarRange: 0}

	queryErr := errors.New("query error")

	rows := sqlmock.NewRows([]string{"id", "displacement", "cylinders", "range"}).
		AddRow(id.String(), 1800, 7, 0)
	mock.ExpectQuery("SELECT *from Engine where id=?").WithArgs(id.String()).WillReturnRows(rows)
	mock.ExpectQuery("SELECT *from Engine where id=?").WithArgs(uuid.Nil).WillReturnError(queryErr)

	testcases := []struct {
		desc   string
		input  uuid.UUID
		output models.Engine
		err    error
	}{
		{"success", engine.EngineID, engine, nil},
		{"failure", uuid.Nil, models.Engine{}, queryErr},
	}
	for i, tc := range testcases {
		resp, err := dbcheck.EngineGetByID(context.TODO(), tc.input.String())

		if resp != tc.output {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, resp, tc.output)
		}

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestEnginestore_EngineCreate function to test createEngine function
func TestEnginestore_EngineCreate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	dbcheck := New(db)

	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	engine := models.Engine{EngineID: id, Displacement: 1600, NoOfCylinder: 4, CarRange: 0}

	queryErr := errors.New("query error")

	mock.ExpectExec("INSERT INTO Engine (id,displacement,cylinders,`range`) VALUES(?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), engine.Displacement, engine.NoOfCylinder, engine.CarRange).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO Engine (id,displacement,cylinders,`range`) VALUES(?,?,?,?)").
		WithArgs(sqlmock.AnyArg(), engine.Displacement, engine.NoOfCylinder, engine.CarRange).
		WillReturnError(queryErr)

	testcases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"success", engine.EngineID, nil},
		{"failure", uuid.Nil, queryErr},
	}

	for i, tc := range testcases {
		_, err := dbcheck.EngineCreate(context.TODO(), &engine)

		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// TestEnginestore_EngineUpdate function to test updateEngine function
func TestEnginestore_EngineUpdate(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Error(err)
	}

	dbcheck := New(db)

	defer func() {
		err = db.Close()
		if err != nil {
			return
		}
	}()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	engine := models.Engine{EngineID: id, Displacement: 1800, NoOfCylinder: 8, CarRange: 1}
	Failed := errors.New("update failed")

	mock.ExpectExec("UPDATE Engine SET displacement=?,cylinders=?,`range`=? WHERE Id=?").
		WithArgs(engine.Displacement, engine.NoOfCylinder, engine.CarRange, engine.EngineID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE Engine SET displacement=?,cylinders=?,`range`=? WHERE Id=?").
		WithArgs(engine.Displacement, engine.NoOfCylinder, engine.CarRange, engine.EngineID).
		WillReturnError(Failed)

	testcases := []struct {
		desc  string
		input models.Engine
		err   error
	}{
		{"success", engine, nil},
		{"failure", engine, Failed},
	}

	for i, tc := range testcases {
		_, err := dbcheck.EngineUpdate(context.TODO(), tc.input.EngineID.String(), tc.input)
		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}

// TestEnginestore_EngineDelete function to test deleteEngine function
func TestEnginestore_EngineDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error(err)
	}

	dbcheck := New(db)

	defer db.Close()

	id, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("cannot generate new id : %v", err)
	}

	deleteErr := errors.New("delete failed")

	mock.ExpectExec("delete from Engine where id=?").WithArgs(id.String()).WillReturnResult(sqlmock.NewResult(
		1, 1))
	mock.ExpectExec("delete  from Engine where id=?").WithArgs(uuid.Nil).WillReturnError(deleteErr)

	cases := []struct {
		desc string
		id   uuid.UUID
		err  error
	}{
		{"Delete success ", id, nil},
		{"Delete failed", uuid.Nil, deleteErr},
	}

	for i, tc := range cases {
		_, err := dbcheck.EngineDelete(context.TODO(), tc.id.String())
		if err != tc.err {
			t.Errorf("\n[TEST %v] Failed \nDesc %v\nGot %v\n Expected %v", i, tc.desc, err, tc.err)
		}
	}
}
